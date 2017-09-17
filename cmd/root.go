package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// RootCmd is the root cobra command for syncrets
var RootCmd = &cobra.Command{
	Use:   "subcommand [src] [dst]",
	Short: "foo is a short command",
	Long:  "Foo is a really short command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TODO: root usage")
	},
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
