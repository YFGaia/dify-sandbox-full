package handlers

import (
	"context"
	"fmt"

	runner_types "github.com/langgenius/dify-sandbox/internal/core/runner/types"
	"github.com/langgenius/dify-sandbox/internal/service"
	"github.com/langgenius/dify-sandbox/internal/types"
	"github.com/mark3labs/mcp-go/mcp"
)

// CodeExecutionHandler handles code execution related MCP tools
type CodeExecutionHandler struct{}

// RunPythonCode executes Python code in the sandbox
func (h *CodeExecutionHandler) RunPythonCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract arguments
	args := request.GetArguments()

	code, ok := args["code"].(string)
	if !ok {
		return mcp.NewToolResultError("Invalid code parameter: must be a string"), nil
	}

	preload := ""
	if p, ok := args["preload"].(string); ok {
		preload = p
	}

	enableNetwork := false
	if en, ok := args["enable_network"].(bool); ok {
		enableNetwork = en
	}

	// Call the service layer directly
	result := service.RunPython3Code(code, preload, &runner_types.RunnerOptions{
		EnableNetwork: enableNetwork,
	})

	// Convert response to MCP result
	return convertDifySandboxResponse(result), nil
}

// RunNodeJSCode executes Node.js code in the sandbox
func (h *CodeExecutionHandler) RunNodeJSCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	// Extract arguments
	args := request.GetArguments()

	code, ok := args["code"].(string)
	if !ok {
		return mcp.NewToolResultError("Invalid code parameter: must be a string"), nil
	}

	preload := ""
	if p, ok := args["preload"].(string); ok {
		preload = p
	}

	enableNetwork := false
	if en, ok := args["enable_network"].(bool); ok {
		enableNetwork = en
	}

	// Call the service layer directly
	result := service.RunNodeJsCode(code, preload, &runner_types.RunnerOptions{
		EnableNetwork: enableNetwork,
	})

	// Convert response to MCP result
	return convertDifySandboxResponse(result), nil
}

// convertDifySandboxResponse converts a DifySandboxResponse to MCP tool result
func convertDifySandboxResponse(response *types.DifySandboxResponse) *mcp.CallToolResult {
	if response.Code != 0 {
		return mcp.NewToolResultError(response.Message)
	}

	// The response.Data is a RunCodeResponse
	if runResult, ok := response.Data.(*service.RunCodeResponse); ok {
		resultText := fmt.Sprintf("Execution completed.\nStdout: %s\nStderr: %s",
			runResult.Stdout, runResult.Stderr)
		return mcp.NewToolResultText(resultText)
	}

	return mcp.NewToolResultText("Execution completed successfully")
}
