package tools

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

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
	if len(args) < 1 {
		return nil, fmt.Errorf("cmd_executor requires 1 arguments, got %d", len(args))
	}
	command, ok := args[0].(string)
	if !ok {
		return nil, fmt.Errorf("command must be a string, got %T", args[0])
	}

	// print the command to be executed
	utils.PrintlnTitle("=>", "CMD:")
	utils.PrintMessage(command)
	// If user input is required, ask for confirmation
	if !waitingUserAllow() {
		return nil, fmt.Errorf("command execution aborted by user")
	}

	// Execute the command
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	cmd := exec.CommandContext(ctx, "powershell.exe", command)
	cmd.Stdin = os.Stdin
	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)
	go exitCall(cancel)
	err := cmd.Run()

	// Capture the output and error
	outStr, errStr := stdout.String(), stderr.String()
	if err != nil && err != context.Canceled {
		utils.PrintlnErrTitle("=>", "FAIL:")
		utils.PrintErr(err.Error())
	}
	utils.PrintlnTitle("=>", "STDOUT:")
	if len(outStr) > 100 {
		utils.PrintMessage(fmt.Sprintf("%s\n...(truncated)", outStr[:100]))
	} else {
		utils.PrintMessage(outStr)
	}
	if errStr != "" {
		utils.PrintlnWarnTitle("=>", "STDERR:")
		if len(errStr) > 100 {
			utils.PrintWarn(fmt.Sprintf("%s\n...(truncated)", errStr[:100]))
		} else {
			utils.PrintWarn(errStr)
		}
	}
	return []string{outStr[:min(100, len(outStr))], errStr[:min(100, len(errStr))]}, nil
}

func waitingUserAllow() bool {
	utils.PrintTitle("<=", fmt.Sprintf("Allow to execute command? [cancel by ctrl+c] (%s/%s): ", color.GreenString("y"), color.RedString("n")))
	var input string
	fmt.Scanln(&input)
	for range 3 {
		if strings.ToLower(input) == "y" || strings.ToLower(input) == "yes" {
			return true
		}
		if strings.ToLower(input) == "n" || strings.ToLower(input) == "no" {
			return false
		}
	}
	return false
}

func exitCall(cancel context.CancelFunc) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	<-sigs
	cancel()
}
