package cmd

import (
	"github.com/spf13/cobra"

	"fmt"
)

var (
	version string
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "....",
	Long:  "....",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("version %s\n", version)
	},
}
