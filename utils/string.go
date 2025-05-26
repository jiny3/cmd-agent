package utils

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintTitle(flag, message string) {
	fmt.Printf("%s %s", color.GreenString(flag), message)
}

func PrintlnTitle(flag, message string) {
	fmt.Printf("%s %s\n", color.GreenString(flag), message)
}

func PrintErrTitle(flag, message string) {
	fmt.Printf("%s %s", color.RedString(flag), message)
}

func PrintlnErrTitle(flag, message string) {
	fmt.Printf("%s %s\n", color.RedString(flag), message)
}

func PrintWarnTitle(flag, message string) {
	fmt.Printf("%s %s", color.YellowString(flag), message)
}

func PrintlnWarnTitle(flag, message string) {
	fmt.Printf("%s %s\n", color.YellowString(flag), message)
}

func PrintMessage(message string) {
	color.Cyan("%s", message)
}

func PrintErr(message string) {
	if len(message) > 100 {
		color.Red("%s\n...(truncated)", message[:100])
	} else {
		color.Red(message)
	}
}

func PrintWarn(message string) {
	if len(message) > 100 {
		color.Yellow("%s\n...(truncated)", message[:100])
	} else {
		color.Yellow(message)
	}
}
