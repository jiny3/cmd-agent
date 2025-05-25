package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"google.golang.org/genai"
)

func cmdExecutorFunctionDeclaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        "cmd_executor",
		Description: "Execute a command in the shell",
		Parameters: &genai.Schema{
			Type: genai.Type("object"),
			Properties: map[string]*genai.Schema{
				"command": {
					Type:        "string",
					Description: "command to execute",
					Example:     "echo \"hello world\"",
				},
				"timeout": {
					Type:        "string",
					Description: "timeout in seconds for the command execution, in case the command loops forever",
					Example:     "10",
				},
			},
			Required: []string{
				"command",
			},
		},
	}
}

func CmdExecutorTool() *genai.Tool {
	declarations := []*genai.FunctionDeclaration{
		cmdExecutorFunctionDeclaration(),
	}
	// Register the cmd_executor function
	for _, declaration := range declarations {
		register(declaration.Name, functionHandler(cmdExecutor))
	}
	return &genai.Tool{
		FunctionDeclarations: declarations,
	}
}

func cmdExecutor(args ...any) (any, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("cmd_executor requires 2 arguments, got %d", len(args))
	}
	command, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("command must be a string, got %T", args[0])
	}
	timeout, ok := args[1].(string)
	if !ok {
		return nil, fmt.Errorf("timeout must be a string, got %T", args[1])
	}
	timeoutInt, err := time.ParseDuration(timeout + "s")
	if err != nil {
		return nil, fmt.Errorf("invalid timeout value: %s, error: %v", timeout, err)
	}
	if timeoutInt <= 0 {
		return nil, fmt.Errorf("timeout must be greater than 0, got %d", timeoutInt)
	}

	// print the command to be executed
	fmt.Println("=> CMD:")
	fmt.Println(command)
	// Execute the command, with timeout (default: 10s)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutInt*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	if err != nil {
		fmt.Println("=> FAIL:")
		fmt.Println(err)
		if errStr != "" {
			fmt.Println("=> STDERR:")
			fmt.Println(errStr)
		}
		return nil, err
	}
	fmt.Println("=> SUCCESS:")
	fmt.Println(outStr)
	if errStr != "" {
		fmt.Println("=> STDERR:")
		fmt.Println(errStr)
	}
	return []string{outStr, errStr}, nil
}
