package cmd

import (
	"github.com/drmdrew/syncrets/core"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(authCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with a system providing secrets",
	Long:  `Authenticate with a system providing secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		core.NewVaultBackend(args)
	},
}
