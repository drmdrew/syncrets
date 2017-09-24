package cmd

import (
	"log"

	"github.com/drmdrew/syncrets/backend"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(authCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with a system providing secrets",
	Long:  `Authenticate with a system providing secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		_, err := backend.NewVaultBackend(viper.GetViper(), args)
		if err != nil {
			log.Fatal(err)
		}
	},
}
