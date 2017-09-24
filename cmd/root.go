package cmd

import (
	"github.com/spf13/cobra"
)

// RootCmd is the root cobra command for syncrets
var RootCmd = &cobra.Command{
	Use:   "subcommand [args] ...",
	Short: "subcommand required such as: auth, list, rm, sync",
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}

// Execute the RootCmd
func Execute() {
	RootCmd.Execute()
}
