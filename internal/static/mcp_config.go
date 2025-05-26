package static

import (
	"os"
	"strconv"
	"strings"

	"github.com/langgenius/dify-sandbox/internal/types"
	"gopkg.in/yaml.v3"
)

var mcpServerConfig types.MCPServerConfig

// InitMCPConfig 初始化 MCP 服务器配置
func InitMCPConfig(path string) error {
	mcpServerConfig = types.MCPServerConfig{}

	// 读取配置文件
	configFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer configFile.Close()

	// 解析配置文件
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&mcpServerConfig)
	if err != nil {
		return err
	}

	// 应用环境变量覆盖
	applyMCPEnvOverrides()

	// 设置默认值
	setMCPDefaults()

	// 转换为兼容格式并设置全局配置
	difySandboxGlobalConfigurations = mcpServerConfig.ConvertToLegacyConfig()

	return nil
}

// applyMCPEnvOverrides 应用环境变量覆盖
func applyMCPEnvOverrides() {
	// MCP 传输配置
	if mode := os.Getenv("MCP_TRANSPORT"); mode != "" {
		mcpServerConfig.MCP.Transport.Mode = mode
	}

	if port := os.Getenv("MCP_HTTP_PORT"); port != "" {
		if p, err := strconv.Atoi(port); err == nil {
			mcpServerConfig.MCP.Transport.HTTPPort = p
		}
	}

	if baseURL := os.Getenv("MCP_BASE_URL"); baseURL != "" {
		mcpServerConfig.MCP.Transport.BaseURL = baseURL
	}

	// 执行环境配置
	if maxWorkers := os.Getenv("MAX_WORKERS"); maxWorkers != "" {
		if w, err := strconv.Atoi(maxWorkers); err == nil {
			mcpServerConfig.Execution.MaxWorkers = w
		}
	}

	if maxRequests := os.Getenv("MAX_REQUESTS"); maxRequests != "" {
		if r, err := strconv.Atoi(maxRequests); err == nil {
			mcpServerConfig.Execution.MaxRequests = r
		}
	}

	if timeout := os.Getenv("WORKER_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			mcpServerConfig.Execution.WorkerTimeout = t
		}
	}

	// Python 配置
	if pythonPath := os.Getenv("PYTHON_PATH"); pythonPath != "" {
		mcpServerConfig.Python.Path = pythonPath
	}

	if pythonLibPath := os.Getenv("PYTHON_LIB_PATH"); pythonLibPath != "" {
		mcpServerConfig.Python.LibPaths = strings.Split(pythonLibPath, ",")
	}

	if depsInterval := os.Getenv("PYTHON_DEPS_UPDATE_INTERVAL"); depsInterval != "" {
		mcpServerConfig.Python.DepsUpdateInterval = depsInterval
	}

	// 安全配置
	if enableNetwork := os.Getenv("ENABLE_NETWORK"); enableNetwork != "" {
		if en, err := strconv.ParseBool(enableNetwork); err == nil {
			mcpServerConfig.Security.EnableNetwork = en
		}
	}

	if enablePreload := os.Getenv("ENABLE_PRELOAD"); enablePreload != "" {
		if ep, err := strconv.ParseBool(enablePreload); err == nil {
			mcpServerConfig.Security.EnablePreload = ep
		}
	}

	if allowedSyscalls := os.Getenv("ALLOWED_SYSCALLS"); allowedSyscalls != "" {
		strs := strings.Split(allowedSyscalls, ",")
		ary := make([]int, len(strs))
		for i := range ary {
			if syscall, err := strconv.Atoi(strs[i]); err == nil {
				ary[i] = syscall
			}
		}
		mcpServerConfig.Security.AllowedSyscalls = ary
	}

	// Node.js security configuration
	if disableNodejsSeccomp := os.Getenv("DISABLE_NODEJS_SECCOMP"); disableNodejsSeccomp != "" {
		if dns, err := strconv.ParseBool(disableNodejsSeccomp); err == nil {
			mcpServerConfig.Security.NodeJS.DisableSeccomp = dns
		}
	}

	// 代理配置
	if socks5Proxy := os.Getenv("SOCKS5_PROXY"); socks5Proxy != "" {
		mcpServerConfig.Proxy.Socks5 = socks5Proxy
	}

	if httpProxy := os.Getenv("HTTP_PROXY"); httpProxy != "" {
		mcpServerConfig.Proxy.HTTP = httpProxy
	}

	if httpsProxy := os.Getenv("HTTPS_PROXY"); httpsProxy != "" {
		mcpServerConfig.Proxy.HTTPS = httpsProxy
	}

	// 日志配置
	if showLog := os.Getenv("MCP_SHOW_LOG"); showLog != "" {
		if sl, err := strconv.ParseBool(showLog); err == nil {
			mcpServerConfig.Logging.ShowLog = sl
		}
	}

	if logLevel := os.Getenv("MCP_LOG_LEVEL"); logLevel != "" {
		mcpServerConfig.Logging.Level = logLevel
	}
}

// setMCPDefaults 设置默认值
func setMCPDefaults() {
	// MCP 服务器默认值
	if mcpServerConfig.MCP.Name == "" {
		mcpServerConfig.MCP.Name = "dify-sandbox"
	}

	if mcpServerConfig.MCP.Version == "" {
		mcpServerConfig.MCP.Version = "1.0.0"
	}

	if mcpServerConfig.MCP.Transport.Mode == "" {
		mcpServerConfig.MCP.Transport.Mode = "stdio"
	}

	if mcpServerConfig.MCP.Transport.HTTPPort == 0 {
		mcpServerConfig.MCP.Transport.HTTPPort = 8080
	}

	if mcpServerConfig.MCP.Transport.BaseURL == "" {
		mcpServerConfig.MCP.Transport.BaseURL = "http://localhost:8080"
	}

	// 执行环境默认值
	if mcpServerConfig.Execution.MaxWorkers == 0 {
		mcpServerConfig.Execution.MaxWorkers = 4
	}

	if mcpServerConfig.Execution.MaxRequests == 0 {
		mcpServerConfig.Execution.MaxRequests = 50
	}

	if mcpServerConfig.Execution.WorkerTimeout == 0 {
		mcpServerConfig.Execution.WorkerTimeout = 60
	}

	// Python 默认值
	if mcpServerConfig.Python.Path == "" {
		mcpServerConfig.Python.Path = "/usr/local/bin/python3"
	}

	if len(mcpServerConfig.Python.LibPaths) == 0 {
		// 使用与原配置系统相同的默认值
		mcpServerConfig.Python.LibPaths = []string{
			"/usr/local/lib/python3.10",
			"/usr/lib/python3.10",
			"/usr/lib/python3",
			"/usr/lib/x86_64-linux-gnu",
			"/usr/lib/aarch64-linux-gnu",
			"/etc/ssl/certs/ca-certificates.crt",
			"/etc/nsswitch.conf",
			"/etc/hosts",
			"/etc/resolv.conf",
			"/run/systemd/resolve/stub-resolv.conf",
			"/run/resolvconf/resolv.conf",
			"/etc/localtime",
			"/usr/share/zoneinfo",
			"/etc/timezone",
			"/usr/local/lib/python3.10/site-packages/pandas",
		}
	}

	if mcpServerConfig.Python.DepsUpdateInterval == "" {
		mcpServerConfig.Python.DepsUpdateInterval = "24h"
	}

	// 日志默认值
	if mcpServerConfig.Logging.Level == "" {
		mcpServerConfig.Logging.Level = "info"
	}
}

// GetMCPServerConfig 获取 MCP 服务器配置
func GetMCPServerConfig() types.MCPServerConfig {
	return mcpServerConfig
}

// GetMCPTransportMode 获取 MCP 传输模式
func GetMCPTransportMode() string {
	return mcpServerConfig.MCP.Transport.Mode
}

// GetMCPHTTPPort 获取 MCP HTTP 端口
func GetMCPHTTPPort() int {
	return mcpServerConfig.MCP.Transport.HTTPPort
}

// GetMCPBaseURL 获取 MCP 基础 URL
func GetMCPBaseURL() string {
	return mcpServerConfig.MCP.Transport.BaseURL
}

// ShouldShowLog 是否显示日志
func ShouldShowLog() bool {
	return mcpServerConfig.Logging.ShowLog
}
