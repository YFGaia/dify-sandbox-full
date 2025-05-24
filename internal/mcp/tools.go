package mcp

import (
	"github.com/langgenius/dify-sandbox/internal/mcp/handlers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterAllTools registers all MCP tools to the server
func RegisterAllTools(s *server.MCPServer) {
	// Initialize handlers
	codeHandler := &handlers.CodeExecutionHandler{}
	depHandler := &handlers.DependencyHandler{}
	healthHandler := &handlers.HealthHandler{}

	// Register run_python_code tool
	runPythonTool := mcp.NewTool("run_python_code",
		mcp.WithDescription("在安全沙箱环境中执行 Python 代码"),
		mcp.WithString("code",
			mcp.Required(),
			mcp.Description("要执行的 Python 代码"),
		),
		mcp.WithString("preload",
			mcp.Description("预加载的代码（可选）"),
		),
		mcp.WithBoolean("enable_network",
			mcp.Description("是否启用网络访问"),
			mcp.DefaultBool(false),
		),
	)
	s.AddTool(runPythonTool, codeHandler.RunPythonCode)

	// Register run_nodejs_code tool
	runNodeJSTool := mcp.NewTool("run_nodejs_code",
		mcp.WithDescription("在安全沙箱环境中执行 Node.js 代码"),
		mcp.WithString("code",
			mcp.Required(),
			mcp.Description("要执行的 Node.js 代码"),
		),
		mcp.WithString("preload",
			mcp.Description("预加载的代码（可选）"),
		),
		mcp.WithBoolean("enable_network",
			mcp.Description("是否启用网络访问"),
			mcp.DefaultBool(false),
		),
	)
	s.AddTool(runNodeJSTool, codeHandler.RunNodeJSCode)

	// Register list_python_dependencies tool
	listDepsTool := mcp.NewTool("list_python_dependencies",
		mcp.WithDescription("获取 Python 环境中已安装的依赖包列表"),
	)
	s.AddTool(listDepsTool, depHandler.ListPythonDependencies)

	// Register refresh_python_dependencies tool
	refreshDepsTool := mcp.NewTool("refresh_python_dependencies",
		mcp.WithDescription("刷新并获取最新的 Python 依赖包列表"),
	)
	s.AddTool(refreshDepsTool, depHandler.RefreshPythonDependencies)

	// Register update_python_dependencies tool
	updateDepsTool := mcp.NewTool("update_python_dependencies",
		mcp.WithDescription("更新 Python 依赖环境"),
	)
	s.AddTool(updateDepsTool, depHandler.UpdatePythonDependencies)

	// Register health_check tool
	healthTool := mcp.NewTool("health_check",
		mcp.WithDescription("检查代码执行环境的健康状态"),
	)
	s.AddTool(healthTool, healthHandler.HealthCheck)
}
