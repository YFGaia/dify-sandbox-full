package handlers

import (
	"context"
	"fmt"

	"github.com/langgenius/dify-sandbox/internal/service"
	"github.com/langgenius/dify-sandbox/internal/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// DependencyHandler handles dependency management related MCP tools
type DependencyHandler struct{}

// ListPythonDependencies lists installed Python packages
func (h *DependencyHandler) ListPythonDependencies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Call the service layer directly
	result := service.ListPython3Dependencies()

	// Convert response to MCP result
	return convertDependencyResponse(result), nil
}

// RefreshPythonDependencies refreshes and lists Python packages
func (h *DependencyHandler) RefreshPythonDependencies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Call the service layer directly
	result := service.RefreshPython3Dependencies()

	// Convert response to MCP result
	return convertDependencyResponse(result), nil
}

// UpdatePythonDependencies updates Python dependency environment
func (h *DependencyHandler) UpdatePythonDependencies(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Call the service layer directly
	result := service.UpdateDependencies()

	// Convert response to MCP result
	if result.Code != 0 {
		return mcp.NewToolResultError(result.Message), nil
	}

	return mcp.NewToolResultText("Python dependencies updated successfully"), nil
}

// convertDependencyResponse converts dependency response to MCP tool result
func convertDependencyResponse(response *types.DifySandboxResponse) *mcp.CallToolResult {
	if response.Code != 0 {
		return mcp.NewToolResultError(response.Message)
	}

	// Handle different response types
	switch data := response.Data.(type) {
	case *service.ListDependenciesResponse:
		resultText := "Installed Python dependencies:\n"
		for _, dep := range data.Dependencies {
			resultText += fmt.Sprintf("- %s: %s\n", dep.Name, dep.Version)
		}
		return mcp.NewToolResultText(resultText)

	case *service.RefreshDependenciesResponse:
		resultText := "Refreshed Python dependencies:\n"
		for _, dep := range data.Dependencies {
			resultText += fmt.Sprintf("- %s: %s\n", dep.Name, dep.Version)
		}
		return mcp.NewToolResultText(resultText)

	default:
		return mcp.NewToolResultText("Operation completed successfully")
	}
}
