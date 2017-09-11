package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"

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
		vault, err := newVault(viper.GetViper(), cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		if err := vault.Authenticate(); err != nil {
			log.Fatalf("Authenication failed: %v", err)
		}
		if !vault.IsValid() {
			log.Fatal("Authentication has failed!")
		}
		log.Print("Authentication was successful")
		vault.Store()
	},
}

func resolveAlias(v *viper.Viper, alias string) *url.URL {
	vkey := fmt.Sprintf("vault.%s.url", alias)
	vurl := v.GetString(vkey)
	log.Printf("Checking for alias: %v", vkey)
	if vurl != "" {
		log.Printf("using alias: %s\n", vurl)
		return parseURL(vurl)
	}
	return nil
}

func newVault(v *viper.Viper, cmd *cobra.Command, args []string) (*backend.Vault, error) {
	if len(args) < 1 {
		return nil, errors.New("source argument is missing")
	}
	u := parseURL(args[0])
	if u == nil {
		return nil, errors.New("cannot parse url")
	}
	hostname := u.Hostname()
	u = resolveAlias(v, hostname)
	if u != nil {
		log.Printf("using alias: %v\n", u)
	}
	return backend.NewVault(v, hostname, u)
}
