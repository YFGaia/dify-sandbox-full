package mcp

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/langgenius/dify-sandbox/internal/core/runner/python"
	"github.com/langgenius/dify-sandbox/internal/static"
	"github.com/langgenius/dify-sandbox/internal/utils/log"
	"github.com/mark3labs/mcp-go/server"
)

// CORS middleware for SSE server
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func initConfig() {
	// 使用 MCP 专用配置文件
	err := static.InitMCPConfig("conf/mcp-config.yaml")
	if err != nil {
		log.Panic("failed to init MCP config: %v", err)
	}
	log.Info("MCP config init success")

	err = static.SetupRunnerDependencies()
	if err != nil {
		log.Error("failed to setup runner dependencies: %v", err)
	}
	log.Info("runner dependencies init success")
}

func initDependencies() {
	log.Info("installing python dependencies...")
	dependencies := static.GetRunnerDependencies()
	err := python.InstallDependencies(dependencies.PythonRequirements)
	if err != nil {
		log.Error("failed to install python dependencies: %v", err)
	}
	log.Info("python dependencies installed")

	log.Info("initializing python dependencies sandbox...")
	err = python.PreparePythonDependenciesEnv()
	if err != nil {
		log.Error("failed to initialize python dependencies sandbox: %v", err)
	}
	log.Info("python dependencies sandbox initialized")

	// start a ticker to update python dependencies to keep the sandbox up-to-date
	go func() {
		updateInterval := static.GetDifySandboxGlobalConfigurations().PythonDepsUpdateInterval
		tickerDuration, err := time.ParseDuration(updateInterval)
		if err != nil {
			log.Error("failed to parse python dependencies update interval, skip periodic updates: %v", err)
			return
		}
		ticker := time.NewTicker(tickerDuration)
		for range ticker.C {
			if err := updatePythonDependencies(dependencies); err != nil {
				log.Error("Failed to update Python dependencies: %v", err)
			}
		}
	}()
}

func updatePythonDependencies(dependencies static.RunnerDependencies) error {
	log.Info("Updating Python dependencies...")
	if err := python.InstallDependencies(dependencies.PythonRequirements); err != nil {
		log.Error("Failed to install Python dependencies: %v", err)
		return err
	}
	if err := python.PreparePythonDependenciesEnv(); err != nil {
		log.Error("Failed to prepare Python dependencies environment: %v", err)
		return err
	}
	log.Info("Python dependencies updated successfully.")
	return nil
}

func StartMCPServer() {
	// 初始化配置和依赖
	initConfig()

	// 获取 MCP 配置
	mcpConfig := static.GetMCPServerConfig()

	// 获取传输模式配置（优先使用环境变量，然后使用配置文件）
	transportMode := getEnvOrDefault("MCP_TRANSPORT", mcpConfig.MCP.Transport.Mode)
	httpPort := getEnvOrDefault("MCP_HTTP_PORT", strconv.Itoa(mcpConfig.MCP.Transport.HTTPPort))

	// 启动信息总是显示，不受配置影响
	log.SetShowLog(true)
	log.Info("Starting Dify-Sandbox MCP Server...")
	log.Info("Transport Mode: %s", transportMode)
	if transportMode == "sse" || transportMode == "http" || transportMode == "streamable-http" {
		log.Info("HTTP Port: %s", httpPort)
		log.Info("Base URL: %s", mcpConfig.MCP.Transport.BaseURL)
	}

	initDependencies()

	// 创建 MCP 服务器
	mcpServer := server.NewMCPServer(
		mcpConfig.MCP.Name,
		mcpConfig.MCP.Version,
		server.WithToolCapabilities(true),
		server.WithResourceCapabilities(false, false),
		server.WithPromptCapabilities(false),
	)

	// 注册工具
	RegisterAllTools(mcpServer)
	log.Info("Registered %d MCP tools", 6) // 我们有6个工具

	switch transportMode {
	case "sse", "http":
		// SSE HTTP 模式 (当前推荐的 HTTP 传输方式)
		log.Info("Starting MCP Server in SSE HTTP mode on port %s", httpPort)

		// BaseURL 应该是客户端可以访问的地址
		// 在容器环境中，这应该是宿主机的地址或域名
		baseURL := getEnvOrDefault("MCP_BASE_URL", mcpConfig.MCP.Transport.BaseURL)

		// 创建 SSE 服务器
		sseServer := server.NewSSEServer(mcpServer,
			server.WithBaseURL(baseURL),
			server.WithSSEEndpoint("/sse"),
			server.WithMessageEndpoint("/message"),
		)

		// 创建带 CORS 支持的 HTTP 服务器
		// 服务器绑定到 0.0.0.0 以接受外部连接
		httpServer := &http.Server{
			Addr:    "0.0.0.0:" + httpPort,
			Handler: corsMiddleware(sseServer),
		}

		// HTTP 模式下根据配置决定是否继续显示日志
		log.SetShowLog(static.ShouldShowLog())

		// 启动 SSE 服务器
		log.Info("✅ MCP SSE Server is now running on 0.0.0.0:%s", httpPort)
		log.Info("📡 SSE Endpoint: %s/sse", baseURL)
		log.Info("📤 Message Endpoint: %s/message", baseURL)
		log.Info("🌐 CORS enabled for all origins")

		if err := httpServer.ListenAndServe(); err != nil {
			log.Panic("Failed to start SSE server: %v", err)
		}

	case "streamable-http":
		// StreamableHTTP 模式 (使用 mcp-go v0.30.0 的完整支持)
		log.Info("Starting MCP Server in StreamableHTTP mode on port %s", httpPort)

		// BaseURL 应该是客户端可以访问的地址
		baseURL := getEnvOrDefault("MCP_BASE_URL", mcpConfig.MCP.Transport.BaseURL)

		// 创建 StreamableHTTP 服务器
		streamableServer := server.NewStreamableHTTPServer(mcpServer,
			server.WithEndpointPath("/mcp"),
			server.WithStateLess(true), // 无状态模式，适合容器部署
		)

		// 创建带 CORS 支持的 HTTP 服务器
		// 服务器绑定到 0.0.0.0 以接受外部连接
		httpServer := &http.Server{
			Addr:    "0.0.0.0:" + httpPort,
			Handler: corsMiddleware(streamableServer),
		}

		// HTTP 模式下根据配置决定是否继续显示日志
		log.SetShowLog(static.ShouldShowLog())

		// 启动 StreamableHTTP 服务器
		log.Info("✅ MCP StreamableHTTP Server is now running on 0.0.0.0:%s", httpPort)
		log.Info("🔗 StreamableHTTP Endpoint: %s/mcp", baseURL)
		log.Info("🌐 CORS enabled for all origins")
		log.Info("📋 Transport: StreamableHTTP (stateless)")

		if err := httpServer.ListenAndServe(); err != nil {
			log.Panic("Failed to start StreamableHTTP server: %v", err)
		}

	default:
		// STDIO 模式（默认）
		log.Info("🔌 Starting MCP Server in STDIO mode...")
		log.Info("✅ MCP Server is now running and waiting for requests...")

		// STDIO 模式下禁用日志输出，避免干扰 JSON-RPC 通信
		log.SetShowLog(false)

		if err := server.ServeStdio(mcpServer); err != nil {
			// 错误时重新启用日志
			log.SetShowLog(true)
			log.Panic("MCP server error: %v", err)
		}
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
