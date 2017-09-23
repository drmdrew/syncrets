package cmd

import (
	"fmt"
	"strings"

	"github.com/drmdrew/syncrets/backend"
	"github.com/drmdrew/syncrets/core"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync secrets from vault",
	Long:  `Sync secrets from vault`,
	Run: func(cmd *cobra.Command, args []string) {
		srcArgs := args[0:1]
		src := core.NewVaultBackend(srcArgs)
		dstArgs := args[1:2]
		dst := core.NewVaultBackend(dstArgs)
		sync(src, dst)
	},
}

func sync(src, dst *backend.Endpoint) {
	var prefixes []string
	path := src.Path
	vault := src.Vault
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
					// this leaf has a secret value...
					// ... now print it
					sep := "/"
					if strings.HasSuffix(prefix, "/") {
						sep = ""
					}
					path := fmt.Sprintf("%s%s%s", prefix, sep, s)
					fmt.Printf("%s\n", path)
					// ... so copy it to dst vault
					value, err := vault.GetClient().Read(path)
					if value != nil {
						fmt.Printf("   -> value.Data['value']: %s\n", value.Data["value"])
						data := value.Data
						bValue, bErr := dst.Vault.GetClient().Write(path, data)
						fmt.Printf("   -> written to destination, path=%s, secret=%v, err=%v\n", path, bValue, bErr)
					} else {
						fmt.Printf("   !! err: %v\n", err)
					}
				}
			}
		}
	}
}
