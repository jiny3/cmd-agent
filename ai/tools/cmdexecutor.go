package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
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
					Description: "timeout in milliseconds for the command execution, in case the command loops forever",
					Example:     "10000ms, 5s, 10s, 1m",
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
	// Validate arguments
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
	timeoutInt, err := time.ParseDuration(timeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout value: %s, error: %v", timeout, err)
	}
	if timeoutInt <= 0 {
		return nil, fmt.Errorf("timeout must be greater than 0, got %d", timeoutInt)
	}

	// print the command to be executed
	fmt.Printf("=> CMD (with %.2f seconds):\n", timeoutInt.Seconds())
	commandSlice := strings.Split(command, " ")
	fmt.Println(command)
	// Execute the command, with timeout (default: 10,000 ms)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutInt)
	defer cancel()
	cmd := exec.CommandContext(ctx, commandSlice[0], commandSlice[1:]...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err = cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	if err != nil {
		fmt.Println("=> FAIL:")
		if err == context.DeadlineExceeded {
			fmt.Println("Command timed out after", timeoutInt)
		} else {
			fmt.Println(err)
		}
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
