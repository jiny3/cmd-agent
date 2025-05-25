package tools

import (
	"bytes"
	"fmt"
	"os/exec"

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
	if len(args) == 0 {
		return nil, fmt.Errorf("arguments are required")
	}
	command, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("command must be a string, got %T", args[0])
	}
	// print the command to be executed
	fmt.Println("=> CMD:")
	fmt.Println(command)
	// Execute the command
	cmd := exec.Command("bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout // 标准输出
	cmd.Stderr = &stderr // 标准错误
	err := cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	if err != nil {
		fmt.Println("=> FAIL:", err)
		if errStr != "" {
			fmt.Println("=> ERROR:")
			fmt.Println(errStr)
		}
		return nil, err
	}
	fmt.Println("=> SUCCESS:")
	fmt.Println(outStr)
	if errStr != "" {
		fmt.Println("=> ERROR:")
		fmt.Println(errStr)
	}
	return []string{outStr, errStr}, nil
}
