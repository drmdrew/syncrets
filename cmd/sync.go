package cmd

import (
	"fmt"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"log"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "....",
	Long:  "....",
	Run:   sync,
}

func sync(cmd *cobra.Command, args []string) {

	client, err := vaultapi.NewClient(vaultapi.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	c := client.Logical()
	secret, err := c.List("secret/")
	if err != nil {
		log.Fatal(err)
	}
	if secret != nil {
		keys := secret.Data["keys"].([]interface{})
		for _, key := range keys {
			fmt.Printf("key: %s\n", key)
		}
	}
	// secret, err := c.Read("secret/foo")
}
