package cmd

import (
	"github.com/go-jarvis/jarvis/version"
	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:     "jarvis",
	Short:   "gerneate jarivs project",
	Version: version.Version,

	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	root.AddCommand(rootCmdNew)
}

func Execute() error {
	return root.Execute()
}
