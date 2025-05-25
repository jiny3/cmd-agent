package tools

import (
	"bytes"
	"fmt"
	"os/exec"

	"google.golang.org/genai"
)

const (
	// cmdExecutorToolName is the name of the command executor tool.
	cmdExecutorToolName = "cmd_executor"
	// cmdExecutorToolDescription is the description of the command executor tool.
	cmdExecutorToolDescription = "Execute a command in the shell"
)

func cmdExecutorToolDeclaration() *genai.FunctionDeclaration {
	return &genai.FunctionDeclaration{
		Name:        cmdExecutorToolName,
		Description: cmdExecutorToolDescription,
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
	return &genai.Tool{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			cmdExecutorToolDeclaration(),
		},
	}
}

func CmdExecutor(command string) (string, string, error) {
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
		return "", errStr, err
	}
	fmt.Println("=> SUCCESS:")
	fmt.Println(outStr)
	if errStr != "" {
		fmt.Println("=> ERROR:")
		fmt.Println(errStr)
	}
	return outStr, errStr, nil
}
