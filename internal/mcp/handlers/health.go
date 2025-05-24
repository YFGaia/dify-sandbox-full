package handlers

import (
	"context"
	"fmt"
	"runtime"

	"github.com/langgenius/dify-sandbox/internal/static"
	"github.com/mark3labs/mcp-go/mcp"
)

// HealthHandler handles health check related MCP tools
type HealthHandler struct{}

// HealthCheck checks the health status of the code execution environment
func (h *HealthHandler) HealthCheck(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config := static.GetDifySandboxGlobalConfigurations()

	healthInfo := fmt.Sprintf(`Health Check Results:
- Service: Dify-Sandbox MCP Server
- Status: Running
- Go Runtime: %s
- Go Version: %s
- Max Workers: %d
- Max Requests: %d
- Worker Timeout: %d seconds
- Python Path: %s
- Network Enabled: %t
- Preload Enabled: %t`,
		runtime.GOOS,
		runtime.Version(),
		config.MaxWorkers,
		config.MaxRequests,
		config.WorkerTimeout,
		config.PythonPath,
		config.EnableNetwork,
		config.EnablePreload,
	)

	return mcp.NewToolResultText(healthInfo), nil
}
