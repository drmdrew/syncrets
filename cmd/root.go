package cmd

import (
	"fmt"
	"log"
	"net/url"
	"os"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
)

// RootCmd is the root cobra command for syncrets
var RootCmd = &cobra.Command{
	Use:   "subcommand [src] [dst]",
	Short: "foo is a short command",
	Long:  "Foo is a really short command",
	Run: func(cmd *cobra.Command, args []string) {
		//		sync(cmd, args)
		fmt.Println("TODO: root usage")
	},
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}

// Execute the RootCmd
func Execute() {

	RootCmd.Execute()
}

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("Cannot parse URL: %v\n", s)
		return nil
	}
	return u
}

// VaultEndpoint client connection to vault
type VaultEndpoint struct {
	client *vaultapi.Client
}

func newVaultEndpoint(u *url.URL, env string) *VaultEndpoint {
	config := vaultapi.DefaultConfig()
	config.Address = fmt.Sprintf("%s://%s", "http", u.Host)
	client, err := vaultapi.NewClient(config)
	if err != nil {
		log.Fatal(err)
	}
	token := os.Getenv(env)
	client.SetToken(token)
	return &VaultEndpoint{client}
}

func (endpoint *VaultEndpoint) list(path string) {
	logical := endpoint.client.Logical()
	secret, err := logical.List(path)
	if err != nil {
		log.Fatal(err)
	}
	if secret != nil {
		keys := secret.Data["keys"].([]interface{})
		for _, key := range keys {
			fmt.Printf("key: %s\n", key)
		}
	}
}

func sync(cmd *cobra.Command, args []string) {

	fmt.Printf("args: %v\n", args)
	if len(args) < 2 {
		log.Fatal("ERROR! source and destination are required!")
	}
	src := parseURL(args[0])
	dst := parseURL(args[1])
	if src == nil || dst == nil {
		log.Fatalf("Error cannot sync from %v to %v\n", src, dst)
	}
	fmt.Printf("source: %s, destination: %s\n", src, dst)

	srcClient := newVaultEndpoint(src, "SYNCRETS_SRC_VAULT_TOKEN")
	dstClient := newVaultEndpoint(dst, "SYNCRETS_DST_VAULT_TOKEN")

	srcClient.list(src.Path)
	dstClient.list(dst.Path)
}
