package cmd

import (
	"fmt"
	"strings"

	"github.com/drmdrew/syncrets/backend"
	"github.com/drmdrew/syncrets/core"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets from vault",
	Long:  `List secrets from vault`,
	Run: func(cmd *cobra.Command, args []string) {
		vault := core.NewVaultBackend(args)
		walk(vault, "secret/")
	},
}

func walk(vault *backend.Vault, path string) {
	var prefixes []string
	prefixes = append(prefixes, path)
	for len(prefixes) > 0 {
		// pop a prefix from the front of the slice
		var prefix string
		prefix, prefixes = prefixes[0], prefixes[1:]
		secret, err := vault.GetClient().List(prefix)
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
