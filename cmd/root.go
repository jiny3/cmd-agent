package cmd

import (
	"os"

	"github.com/jiny3/cmd-agent/ai/client"
	"github.com/jiny3/cmd-agent/ai/tools"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "x <a prompt to generate what you want>",
	Short: "x is a command line tool for AI agents",
	Long:  `an ai agent that uses natural language to execute commands. `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		prompt := args[0]
		resp, err := client.GenerateContent(prompt, tools.CmdExecutorTool())
		if err != nil {
			logrus.Fatal(err)
		}
		logrus.Debug(resp)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
