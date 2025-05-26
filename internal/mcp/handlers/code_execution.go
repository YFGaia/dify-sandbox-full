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

// CodeExecutionHandler handles code execution related MCP tools
type CodeExecutionHandler struct{}

// CodeExecutionResult represents the structured response for code execution
type CodeExecutionResult struct {
	Success   bool   `json:"success"`
	Language  string `json:"language"`
	Stdout    string `json:"stdout"`
	Stderr    string `json:"stderr"`
	Error     string `json:"error,omitempty"`
	Timestamp string `json:"timestamp"`
}

// runCodeWithContext executes code with context timeout control
func runCodeWithContext(ctx context.Context, language string, code string, preload string, options *runner_types.RunnerOptions) (*types.DifySandboxResponse, error) {
	// Create a channel to receive the result
	resultChan := make(chan *types.DifySandboxResponse, 1)

	// Run the code execution in a goroutine
	go func() {
		var result *types.DifySandboxResponse
		switch language {
		case "python":
			result = service.RunPython3Code(code, preload, options)
		case "nodejs":
			result = service.RunNodeJsCode(code, preload, options)
		default:
			result = types.ErrorResponse(-400, "unsupported language")
		}

		select {
		case resultChan <- result:
		case <-ctx.Done():
			// Context cancelled, don't send result
		}
	}()

	// Wait for either result or context cancellation
	select {
	case result := <-resultChan:
		return result, nil
	case <-ctx.Done():
		return types.ErrorResponse(-408, "Request timed out"), ctx.Err()
	}
}

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

	// Call the service layer with context control
	result, err := runCodeWithContext(ctx, "python", code, preload, &runner_types.RunnerOptions{
		EnableNetwork: enableNetwork,
	})

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Execution error: %v", err)), nil
	}

	// Convert response to MCP result
	return convertDifySandboxResponse(result, "python"), nil
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

	// Call the service layer with context control
	result, err := runCodeWithContext(ctx, "nodejs", code, preload, &runner_types.RunnerOptions{
		EnableNetwork: enableNetwork,
	})

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Execution error: %v", err)), nil
	}

	// Convert response to MCP result
	return convertDifySandboxResponse(result, "nodejs"), nil
}

// convertDifySandboxResponse converts a DifySandboxResponse to MCP tool result with structured JSON
func convertDifySandboxResponse(response *types.DifySandboxResponse, language string) *mcp.CallToolResult {
	currentTime := time.Now().Format(time.RFC3339)

	if response.Code != 0 {
		// Create error response
		errorResult := CodeExecutionResult{
			Success:   false,
			Language:  language,
			Stdout:    "",
			Stderr:    "",
			Error:     response.Message,
			Timestamp: currentTime,
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

	// Handle successful response
	if runResult, ok := response.Data.(*service.RunCodeResponse); ok {
		successResult := CodeExecutionResult{
			Success:   true,
			Language:  language,
			Stdout:    runResult.Stdout,
			Stderr:    runResult.Stderr,
			Error:     "",
			Timestamp: currentTime,
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

	// Fallback for unknown response type
	fallbackResult := CodeExecutionResult{
		Success:   true,
		Language:  language,
		Stdout:    "",
		Stderr:    "",
		Error:     "",
		Timestamp: currentTime,
	}

	jsonData, err := json.Marshal(fallbackResult)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to serialize fallback response: %v", err))
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
