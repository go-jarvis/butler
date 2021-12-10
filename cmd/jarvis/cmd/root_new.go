package cmd

import (
	"github.com/go-jarvis/cobrautils"
	"github.com/go-jarvis/jarvis/pkg/jarvis"
	"github.com/spf13/cobra"
)

var rootCmdNew = &cobra.Command{
	Use:   "new",
	Short: "create a new project",
	Run: func(cmd *cobra.Command, args []string) {
		jarvis.Project.CreateProject()
	},
}

func init() {
	cobrautils.BindFlags(rootCmdNew, jarvis.Project)
}
