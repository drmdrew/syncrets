package cmd

import (
	"fmt"
	"log"

	vaultapi "github.com/hashicorp/vault/api"
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
		auth(cmd, args)
		//fmt.Println("TODO: authenticate functionality.")
	},
}

func auth(cmd *cobra.Command, args []string) {

	fmt.Printf("args: %v\n", args)
	if len(args) < 1 {
		log.Fatal("ERROR! source required!")
	}
	src := parseUrl(args[0])
	fmt.Printf("source: %s\n", src)

	//	srcClient := newVaultEndpoint(src, "SYNCRETS_SRC_VAULT_TOKEN")
	config := vaultapi.DefaultConfig()
	config.Address = fmt.Sprintf("%s://%s", "http", src.Host)
	client, err := vaultapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	auth := client.Auth()
	fmt.Printf("auth: %v\n", auth)

}
