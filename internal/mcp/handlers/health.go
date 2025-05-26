package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/langgenius/dify-sandbox/internal/static"
	"github.com/mark3labs/mcp-go/mcp"
)

// HealthHandler handles health check related MCP tools
type HealthHandler struct{}

// HealthCheckResult represents the structured response for health check
type HealthCheckResult struct {
	Success   bool        `json:"success"`
	Service   string      `json:"service"`
	Status    string      `json:"status"`
	Runtime   RuntimeInfo `json:"runtime"`
	Config    ConfigInfo  `json:"config"`
	Timestamp string      `json:"timestamp"`
}

// RuntimeInfo contains runtime information
type RuntimeInfo struct {
	OS        string `json:"os"`
	Version   string `json:"version"`
	GoVersion string `json:"go_version"`
}

// ConfigInfo contains configuration information
type ConfigInfo struct {
	MaxWorkers     int    `json:"max_workers"`
	MaxRequests    int    `json:"max_requests"`
	WorkerTimeout  int    `json:"worker_timeout"`
	PythonPath     string `json:"python_path"`
	NetworkEnabled bool   `json:"network_enabled"`
	PreloadEnabled bool   `json:"preload_enabled"`
}

// HealthCheck checks the health status of the code execution environment
func (h *HealthHandler) HealthCheck(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	config := static.GetDifySandboxGlobalConfigurations()
	currentTime := time.Now().Format(time.RFC3339)

	healthResult := HealthCheckResult{
		Success: true,
		Service: "Dify-Sandbox MCP Server",
		Status:  "Running",
		Runtime: RuntimeInfo{
			OS:        runtime.GOOS,
			Version:   runtime.Version(),
			GoVersion: runtime.Version(),
		},
		Config: ConfigInfo{
			MaxWorkers:     config.MaxWorkers,
			MaxRequests:    config.MaxRequests,
			WorkerTimeout:  config.WorkerTimeout,
			PythonPath:     config.PythonPath,
			NetworkEnabled: config.EnableNetwork,
			PreloadEnabled: config.EnablePreload,
		},
		Timestamp: currentTime,
	}

	jsonData, err := json.Marshal(healthResult)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize health check response: %v", err)), nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(jsonData),
			},
		},
		IsError: false,
	}, nil
}
