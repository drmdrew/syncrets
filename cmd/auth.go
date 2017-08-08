package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VaultReader is just the Read portion of the Vault client API
type VaultReader interface {
	Read(path string) (*vaultapi.Secret, error)
}

// VaultClient is a composite API of all the Vault client APIs as interfaces
type VaultClient interface {
	VaultReader
}

type authenticator struct {
	url    *url.URL
	token  string
	viper  *viper.Viper
	reader VaultReader //*vaultapi.Logical
}

type vaultClient struct {
	logical *vaultapi.Logical
}

func (vc vaultClient) Read(path string) (*vaultapi.Secret, error) {
	return vc.logical.Read(path)
}

func init() {
	RootCmd.AddCommand(authCmd)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with a system providing secrets",
	Long:  `Authenticate with a system providing secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		auth, err := newAuthenticator(viper.GetViper(), cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		if !auth.isValid() {
			log.Fatal("Authentication has failed!")
		}
		log.Print("Authentication was successful")
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

func newAuthenticator(v *viper.Viper, cmd *cobra.Command, args []string) (*authenticator, error) {
	auth := &authenticator{}
	auth.viper = v
	if len(args) < 1 {
		return nil, errors.New("source argument is missing")
	}
	auth.url = parseURL(args[0])
	if auth.url == nil {
		return nil, errors.New("cannot parse url")
	}
	log.Printf("source: %s\n", auth.url.Hostname())
	url := resolveAlias(v, auth.url.Hostname())
	if url != nil {
		log.Printf("using alias: %v\n", url)
		auth.url = url
	}
	client, err := newVaultClient(auth.url)
	if err != nil {
		return nil, err
	}
	auth.reader = *client
	return auth, nil
}

func newVaultClient(src *url.URL) (*vaultClient, error) {
	config := vaultapi.DefaultConfig()
	config.Address = fmt.Sprintf("%s://%s", "http", src.Host)
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	// TODO: clean this up. prompt for token if not set (for now)
	var token string
	fmt.Printf("token: ")
	fmt.Scanf("%s", &token)
	client.SetToken(token)
	vc := &vaultClient{client.Logical()}
	return vc, nil
}

func (auth *authenticator) isValid() bool {
	// use lookup-self to verify token is valid
	secret, err := auth.reader.Read("auth/token/lookup-self")
	if err != nil {
		log.Printf("lookup-self failed: %v\n", err)
		return false
	}
	log.Printf("token id: %v\n", secret.Data["id"])
	return true
}
