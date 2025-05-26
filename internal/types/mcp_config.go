package types

// MCPServerConfig MCP 服务器专用配置结构
type MCPServerConfig struct {
	MCP struct {
		Name      string `yaml:"name"`
		Version   string `yaml:"version"`
		Transport struct {
			Mode     string `yaml:"mode"`
			HTTPPort int    `yaml:"http_port"`
			BaseURL  string `yaml:"base_url"`
		} `yaml:"transport"`
	} `yaml:"mcp"`

	Execution struct {
		MaxWorkers    int `yaml:"max_workers"`
		MaxRequests   int `yaml:"max_requests"`
		WorkerTimeout int `yaml:"worker_timeout"`
	} `yaml:"execution"`

	Python struct {
		Path               string   `yaml:"path"`
		LibPaths           []string `yaml:"lib_paths"`
		DepsUpdateInterval string   `yaml:"deps_update_interval"`
	} `yaml:"python"`

	Security struct {
		EnableNetwork   bool  `yaml:"enable_network"`
		EnablePreload   bool  `yaml:"enable_preload"`
		AllowedSyscalls []int `yaml:"allowed_syscalls"`
		NodeJS          struct {
			DisableSeccomp bool `yaml:"disable_seccomp"`
		} `yaml:"nodejs"`
	} `yaml:"security"`

	Proxy struct {
		Socks5 string `yaml:"socks5"`
		HTTP   string `yaml:"http"`
		HTTPS  string `yaml:"https"`
	} `yaml:"proxy"`

	Logging struct {
		ShowLog bool   `yaml:"show_log"`
		Level   string `yaml:"level"`
	} `yaml:"logging"`
}

// ConvertToLegacyConfig 将 MCP 配置转换为兼容原有系统的配置格式
func (mcp *MCPServerConfig) ConvertToLegacyConfig() DifySandboxGlobalConfigurations {
	return DifySandboxGlobalConfigurations{
		App: struct {
			Port  int    `yaml:"port"`
			Debug bool   `yaml:"debug"`
			Key   string `yaml:"key"`
		}{
			Port:  mcp.MCP.Transport.HTTPPort, // 虽然 MCP 不使用，但保持兼容性
			Debug: mcp.Logging.Level == "debug",
			Key:   "mcp-server", // MCP 不需要 API key，但保持兼容性
		},
		MaxWorkers:               mcp.Execution.MaxWorkers,
		MaxRequests:              mcp.Execution.MaxRequests,
		WorkerTimeout:            mcp.Execution.WorkerTimeout,
		PythonPath:               mcp.Python.Path,
		PythonLibPaths:           mcp.Python.LibPaths,
		PythonDepsUpdateInterval: mcp.Python.DepsUpdateInterval,
		NodejsPath:               "/usr/local/bin/node", // 默认 Node.js 路径
		EnableNetwork:            mcp.Security.EnableNetwork,
		EnablePreload:            mcp.Security.EnablePreload,
		AllowedSyscalls:          mcp.Security.AllowedSyscalls,
		Proxy: struct {
			Socks5 string `yaml:"socks5"`
			Https  string `yaml:"https"`
			Http   string `yaml:"http"`
		}{
			Socks5: mcp.Proxy.Socks5,
			Https:  mcp.Proxy.HTTPS,
			Http:   mcp.Proxy.HTTP,
		},
	}
}
