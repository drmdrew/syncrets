package cmd

import (
	"log"
	"net/url"

	"github.com/drmdrew/syncrets/backend"

	"github.com/drmdrew/syncrets/core"
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
		var hostname string
		var url *url.URL
		var vault *backend.Vault
		var err error
		if hostname, url, err = core.ResolveArgs(viper.GetViper(), args); err != nil {
			log.Fatal(err)
		}
		if vault, err = backend.NewVault(viper.GetViper(), hostname, url); err != nil {
			log.Fatal(err)
		}
		if err = vault.Authenticate(); err != nil {
			log.Fatalf("Authenication failed: %v", err)
		}
		if !vault.IsValid() {
			log.Fatal("Authentication has failed!")
		}
		log.Print("Authentication was successful")
		vault.Store()
	},
}
