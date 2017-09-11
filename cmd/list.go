package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets from vault",
	Long:  `List secrets from vault`,
	Run: func(cmd *cobra.Command, args []string) {
		auth, err := newAuthenticator(viper.GetViper(), cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		if err := auth.authenticate(); err != nil {
			log.Fatalf("Authenication failed: %v", err)
		}
		if !auth.isValid() {
			log.Fatal("Authentication has failed!")
		}
		log.Print("Authentication was successful")
		auth.store()
		walk(auth, "secret/")
	},
}

func walk(auth *authenticator, path string) {
	var prefixes []string
	prefixes = append(prefixes, path)
	for len(prefixes) > 0 {
		// pop a prefix from the front of the slice
		var prefix string
		prefix, prefixes = prefixes[0], prefixes[1:]
		secret, err := auth.client.List(prefix)
		if err != nil {
			continue
		}
		if secret != nil {
			for _, val := range secret.Data["keys"].([]interface{}) {
				s := val.(string)
				if strings.HasSuffix(s, "/") {
					// push a new prefix at the end of the slice
					prefixes = append(prefixes, prefix+s)
				} else {
					fmt.Printf("%s%s\n", prefix, s)
				}
			}
		}
	}
}