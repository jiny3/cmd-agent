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
	fmt.Println(color.CyanString(message))
}

func PrintErr(message string) {
	fmt.Println(color.RedString(message))
}

func PrintWarn(message string) {
	fmt.Println(color.YellowString(message))
}
