package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "foo",
	Short: "foo is a short command",
	Long:  "Foo is a really short command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("RootCmd.Run was called\n")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}

func Execute() {

	RootCmd.Execute()
}
