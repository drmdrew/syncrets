package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	vaultapi "github.com/hashicorp/vault/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// VaultReader is just the Read portion of the Vault client API
type VaultReader interface {
	Read(path string) (*vaultapi.Secret, error)
	Authenticate(token string) error
}

// VaultClient is a composite API of all the Vault client APIs as interfaces
type VaultClient interface {
	VaultReader
}

type authenticator struct {
	hostname string
	url      *url.URL
	token    string
	viper    *viper.Viper
	client   VaultClient
}

type vaultClient struct {
	client *vaultapi.Client
}

func (vc *vaultClient) Read(path string) (*vaultapi.Secret, error) {
	return vc.client.Logical().Read(path)
}

func (vc *vaultClient) Authenticate(token string) error {
	if token == "" {
		// TODO: find a better solution to prompt for token
		fmt.Printf("token: ")
		fmt.Scanf("%s", &token)
	}
	vc.client.SetToken(token)
	return nil
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
	auth.hostname = auth.url.Hostname()
	url := resolveAlias(v, auth.url.Hostname())
	if url != nil {
		log.Printf("using alias: %v\n", url)
		auth.url = url
	}
	client, err := newVaultClient(auth.url)
	if err != nil {
		return nil, err
	}
	auth.client = client
	return auth, nil
}

func newVaultClient(src *url.URL) (*vaultClient, error) {
	config := vaultapi.DefaultConfig()
	config.Address = fmt.Sprintf("%s://%s", "http", src.Host)
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	vc := &vaultClient{client}
	return vc, nil
}

func (auth *authenticator) isValid() bool {
	// load vault token from token-file if one is present
	vkey := fmt.Sprintf("vault.%s.token-file", auth.hostname)
	tokenFile := auth.viper.GetString(vkey)
	token, err := ioutil.ReadFile(tokenFile)
	log.Printf("%v is configured: %v\n", vkey, tokenFile)
	if token != nil {
		trimmedToken := strings.TrimSpace(string(token))
		log.Printf("token is %s\n", trimmedToken)
		auth.client.Authenticate(trimmedToken)
	} else {
		// make sure the vault client has authenticated
		if err := auth.client.Authenticate(""); err != nil {
			log.Printf("authenticate failed: %v\n", err)
			return false
		}
	}

	// use lookup-self to verify token is valid
	secret, err := auth.client.Read("auth/token/lookup-self")
	if err != nil {
		log.Printf("lookup-self failed: %v\n", err)
		return false
	}
	log.Printf("token id: %v\n", secret.Data["id"])
	return true
}
