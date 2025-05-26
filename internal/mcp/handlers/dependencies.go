package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	runner_types "github.com/langgenius/dify-sandbox/internal/core/runner/types"
	"github.com/langgenius/dify-sandbox/internal/service"
	"github.com/langgenius/dify-sandbox/internal/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// DependencyHandler handles dependency management related MCP tools
type DependencyHandler struct{}

// DependencyOperationResult represents the structured response for dependency operations
type DependencyOperationResult struct {
	Success      bool                      `json:"success"`
	Operation    string                    `json:"operation"`
	Dependencies []runner_types.Dependency `json:"dependencies,omitempty"`
	Message      string                    `json:"message,omitempty"`
	Error        string                    `json:"error,omitempty"`
	Timestamp    string                    `json:"timestamp"`
}

// ListPythonDependencies lists installed Python packages
func (h *DependencyHandler) ListPythonDependencies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Call the service layer directly
	result := service.ListPython3Dependencies()

	// Convert response to MCP result
	return convertDependencyResponse(result, "list"), nil
}

// RefreshPythonDependencies refreshes and lists Python packages
func (h *DependencyHandler) RefreshPythonDependencies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Call the service layer directly
	result := service.RefreshPython3Dependencies()

	// Convert response to MCP result
	return convertDependencyResponse(result, "refresh"), nil
}

// UpdatePythonDependencies updates Python dependency environment
func (h *DependencyHandler) UpdatePythonDependencies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Call the service layer directly
	result := service.UpdateDependencies()

	// Convert response to MCP result
	return convertDependencyResponse(result, "update"), nil
}

// convertDependencyResponse converts dependency response to MCP tool result with structured JSON
func convertDependencyResponse(response *types.DifySandboxResponse, operation string) *mcp.CallToolResult {
	currentTime := time.Now().Format(time.RFC3339)

	if response.Code != 0 {
		// Create error response
		errorResult := DependencyOperationResult{
			Success:      false,
			Operation:    operation,
			Dependencies: nil,
			Message:      "",
			Error:        response.Message,
			Timestamp:    currentTime,
		}

		jsonData, err := json.Marshal(errorResult)
		if err != nil {
			return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize error response: %v", err))
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(jsonData),
				},
			},
			IsError: true,
		}
	}

	// Handle different response types
	var successResult DependencyOperationResult

	switch data := response.Data.(type) {
	case *service.ListDependenciesResponse:
		successResult = DependencyOperationResult{
			Success:      true,
			Operation:    operation,
			Dependencies: data.Dependencies,
			Message:      fmt.Sprintf("Found %d installed Python dependencies", len(data.Dependencies)),
			Error:        "",
			Timestamp:    currentTime,
		}

	case *service.RefreshDependenciesResponse:
		successResult = DependencyOperationResult{
			Success:      true,
			Operation:    operation,
			Dependencies: data.Dependencies,
			Message:      fmt.Sprintf("Refreshed %d Python dependencies", len(data.Dependencies)),
			Error:        "",
			Timestamp:    currentTime,
		}

	case *service.UpdateDependenciesResponse:
		successResult = DependencyOperationResult{
			Success:      true,
			Operation:    operation,
			Dependencies: nil,
			Message:      "Python dependencies updated successfully",
			Error:        "",
			Timestamp:    currentTime,
		}

	default:
		successResult = DependencyOperationResult{
			Success:      true,
			Operation:    operation,
			Dependencies: nil,
			Message:      "Operation completed successfully",
			Error:        "",
			Timestamp:    currentTime,
		}
	}

	jsonData, err := json.Marshal(successResult)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize success response: %v", err))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: string(jsonData),
			},
		},
		IsError: false,
	}
}
