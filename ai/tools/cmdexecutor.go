package tools

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jiny3/cmd-agent/utils"
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
	utils.PrintlnTitle("=>", fmt.Sprintf("CMD (with %.2f seconds):", timeoutInt.Seconds()))
	utils.PrintMessage(command)
	// If user input is required, ask for confirmation
	if !waitingUserInput() {
		return nil, fmt.Errorf("command execution aborted by user")
	}

	// Execute the command, with timeout (default: 10,000 ms)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutInt)
	defer cancel()
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// Run the command
	err = cmd.Run()
	outStr, errStr := stdout.String(), stderr.String()
	if err != nil {
		utils.PrintlnErrTitle("=>", "FAIL:")
		if err == context.DeadlineExceeded {
			utils.PrintErr(fmt.Sprintf("Command timed out after %.2f", timeoutInt.Seconds()))
		} else {
			utils.PrintErr(err.Error())
		}
		if errStr != "" {
			utils.PrintlnWarnTitle("=>", "STDERR:")
			utils.PrintWarn(errStr)
		}
		return nil, err
	}
	utils.PrintlnTitle("=>", "SUCCESS:")
	utils.PrintMessage(outStr)
	if errStr != "" {
		utils.PrintlnWarnTitle("=>", "STDERR:")
		utils.PrintWarn(errStr)
	}
	return []string{outStr, errStr}, nil
}

func waitingUserInput() bool {
	utils.PrintTitle("<=", fmt.Sprintf("Allow to execute command? (%s/%s): ", color.GreenString("Y"), color.RedString("N")))
	var input string
	fmt.Scanln(&input)
	if strings.ToLower(input) == "y" || strings.ToLower(input) == "yes" {
		return true
	}
	return false
}
