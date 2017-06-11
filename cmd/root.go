package cmd

import (
	"fmt"
	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"log"
	"net/url"
	"os"
)

var (
	version string
)

var RootCmd = &cobra.Command{
	Use:   "foo [src] [dst]",
	Short: "foo is a short command",
	Long:  "Foo is a really short command",
	Run: func(cmd *cobra.Command, args []string) {
		sync(cmd, args)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}

func Execute() {

	RootCmd.Execute()
}

func parseUrl(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		log.Fatal(err)
	}
	return u
}

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
	src := parseUrl(args[0])
	dst := parseUrl(args[1])
	fmt.Printf("source: %s, destination: %s\n", src, dst)

	srcClient := newVaultEndpoint(src, "SYNCRETS_SRC_VAULT_TOKEN")
	dstClient := newVaultEndpoint(dst, "SYNCRETS_DST_VAULT_TOKEN")

	srcClient.list(src.Path)
	dstClient.list(dst.Path)
}

func printVersion() {
	fmt.Printf("version %s\n", version)
}
