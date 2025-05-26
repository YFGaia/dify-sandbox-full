package nodejs

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/langgenius/dify-sandbox/internal/core/runner"
	"github.com/langgenius/dify-sandbox/internal/core/runner/types"
	"github.com/langgenius/dify-sandbox/internal/static"
)

type NodeJsRunner struct {
	runner.TempDirRunner
}

//go:embed prescript.js
var nodejs_sandbox_fs []byte

var (
	REQUIRED_FS = []string{
		path.Join(LIB_PATH, PROJECT_NAME, "node_temp"),
		path.Join(LIB_PATH, LIB_NAME),
		"/etc/ssl/certs/ca-certificates.crt",
		"/etc/nsswitch.conf",
		"/etc/resolv.conf",
		"/run/systemd/resolve/stub-resolv.conf",
		"/etc/hosts",
	}
)

func (p *NodeJsRunner) Run(
	code string,
	timeout time.Duration,
	stdin []byte,
	preload string,
	options *types.RunnerOptions,
) (chan []byte, chan []byte, chan bool, error) {
	configuration := static.GetDifySandboxGlobalConfigurations()

	// Early environment check to fail fast
	if err := checkNodeJsEnvironment(); err != nil {
		return nil, nil, nil, fmt.Errorf("Node.js environment check failed: %w", err)
	}

	// capture the output
	output_handler := runner.NewOutputCaptureRunner()
	output_handler.SetTimeout(timeout)

	err := p.WithTempDir("/", REQUIRED_FS, func(root_path string) error {
		output_handler.SetAfterExitHook(func() {
			os.RemoveAll(root_path)
			os.Remove(root_path)
		})

		// initialize the environment
		script_path, err := p.InitializeEnvironment(code, preload, root_path)
		if err != nil {
			return fmt.Errorf("failed to initialize Node.js environment: %w", err)
		}

		// Verify Node.js path exists
		nodejsPath := static.GetDifySandboxGlobalConfigurations().NodejsPath
		if _, err := os.Stat(nodejsPath); os.IsNotExist(err) {
			return fmt.Errorf("Node.js binary not found at %s", nodejsPath)
		}

		// Log the script content for debugging
		if scriptContent, err := os.ReadFile(script_path); err == nil {
			fmt.Printf("DEBUG: Node.js script length: %d bytes\n", len(scriptContent))
		}

		// create a new process
		cmd := exec.Command(
			nodejsPath,
			script_path,
			strconv.Itoa(static.SANDBOX_USER_UID),
			strconv.Itoa(static.SANDBOX_GROUP_ID),
			options.Json(),
		)
		cmd.Env = []string{}

		// Add initialization timeout environment variable
		cmd.Env = append(cmd.Env, "NODEJS_INIT_TIMEOUT=5000") // 5 seconds - much shorter

		// Add option to completely disable Node.js seccomp
		// Check both environment variable and MCP config
		mcpConfig := static.GetMCPServerConfig()
		if os.Getenv("DISABLE_NODEJS_SECCOMP") == "true" || mcpConfig.Security.NodeJS.DisableSeccomp {
			cmd.Env = append(cmd.Env, "DISABLE_NODEJS_SECCOMP=true")
			fmt.Printf("INFO: Node.js seccomp completely disabled\n")
		}

		// Add debug option to skip seccomp (for debugging purposes)
		// This can be controlled via environment variable
		if os.Getenv("DEBUG_SKIP_SECCOMP") == "true" {
			cmd.Env = append(cmd.Env, "SKIP_SECCOMP=true")
			fmt.Printf("DEBUG: Skipping seccomp for Node.js execution\n")
		}

		if len(configuration.AllowedSyscalls) > 0 {
			cmd.Env = append(
				cmd.Env,
				fmt.Sprintf("ALLOWED_SYSCALLS=%s", strings.Trim(
					strings.Join(strings.Fields(fmt.Sprint(configuration.AllowedSyscalls)), ","), "[]",
				)),
			)
		}

		// Add NODE_PATH to help with module resolution
		nodeModulesPath := path.Join(root_path, LIB_PATH, PROJECT_NAME, "node_temp/node_modules")
		cmd.Env = append(cmd.Env, fmt.Sprintf("NODE_PATH=%s", nodeModulesPath))

		fmt.Printf("DEBUG: Starting Node.js process with timeout: %v\n", timeout)

		// capture the output
		err = output_handler.CaptureOutput(cmd)
		if err != nil {
			return fmt.Errorf("failed to capture Node.js output: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, nil, err
	}

	return output_handler.GetStdout(), output_handler.GetStderr(), output_handler.GetDone(), nil
}

func (p *NodeJsRunner) InitializeEnvironment(code string, preload string, root_path string) (string, error) {
	if !checkLibAvaliable() {
		releaseLibBinary()
	}

	node_sandbox_file := string(nodejs_sandbox_fs)
	if preload != "" {
		node_sandbox_file = fmt.Sprintf("%s\n%s", preload, node_sandbox_file)
	}

	// join nodejs_sandbox_fs and code
	// encode code with base64
	code = base64.StdEncoding.EncodeToString([]byte(code))
	// FIXE: redeclared function causes code injection
	evalCode := fmt.Sprintf("eval(Buffer.from('%s', 'base64').toString('utf-8'))", code)
	code = node_sandbox_file + evalCode

	// override root_path/tmp/sandbox-nodejs-project/prescript.js
	script_path := path.Join(root_path, LIB_PATH, PROJECT_NAME, "node_temp/node_temp/test.js")
	err := os.WriteFile(script_path, []byte(code), 0755)
	if err != nil {
		return "", err
	}

	return script_path, nil
}
