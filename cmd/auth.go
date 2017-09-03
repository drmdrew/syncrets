package cmd

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/drmdrew/syncrets/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type authenticator struct {
	hostname string
	url      *url.URL
	token    string
	viper    *viper.Viper
	client   vault.ClientAPI
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
		if err := auth.authenticate(); err != nil {
			log.Fatalf("Authenication failed: %v", err)
		}
		if !auth.isValid() {
			log.Fatal("Authentication has failed!")
		}
		log.Print("Authentication was successful")
		auth.store()
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
	client, err := vault.NewClient(auth.url)
	if err != nil {
		return nil, err
	}
	auth.client = client
	return auth, nil
}

// prompt the user for information
func (auth *authenticator) prompt(prompt string) string {
	var token string
	fmt.Printf(prompt)
	fmt.Scanf("%s", &token)
	return token
}

func (auth *authenticator) authenticate() error {
	auth.load()
	if auth.isValid() {
		return nil
	}
	// re-authenticate if loaded token is invalid
	vkey := fmt.Sprintf("vault.%s.auth.method", auth.hostname)
	method := auth.viper.GetString(vkey)
	switch method {
	case "token":
		auth.tokenAuth()
	case "userpass":
		auth.userpassAuth()
	default:
		return fmt.Errorf("No valid auth.method configured for %s", auth.hostname)
	}
	return nil
}

func (auth *authenticator) tokenAuth() {
	token := auth.prompt("token: ")
	auth.client.SetToken(token)
}

func (auth *authenticator) userpassAuth() {
	vkey := fmt.Sprintf("vault.%s.auth.username", auth.hostname)
	username := auth.viper.GetString(vkey)
	password := auth.prompt("password: ")
	err := auth.client.UserpassLogin(username, password)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
}

func (auth *authenticator) isValid() bool {
	// use lookup-self to verify token is valid
	secret, err := auth.client.Read("auth/token/lookup-self")
	if err != nil {
		log.Printf("lookup-self failed: %v\n", err)
		return false
	}
	log.Printf("token id: %v\n", secret.Data["id"])
	return true
}

func (auth *authenticator) load() (string, error) {
	// load vault token from token.file if one is present
	vkey := fmt.Sprintf("vault.%s.token.file", auth.hostname)
	tokenFile := auth.viper.GetString(vkey)
	tokenBytes, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return "", err
	}
	log.Printf("%v is configured: %v\n", vkey, tokenFile)
	var token string
	if tokenBytes != nil {
		token = strings.TrimSpace(string(tokenBytes))
		log.Printf("token is %s\n", token)
	} else {
		return "", fmt.Errorf("Unable to read token from %s", tokenFile)
	}
	auth.client.SetToken(token)
	return token, nil
}

func (auth *authenticator) store() {
	// store the vault token in token.file if one is present
	vkey := fmt.Sprintf("vault.%s.token.file", auth.hostname)
	tokenFile := auth.viper.GetString(vkey)
	token := auth.client.GetToken()
	err := ioutil.WriteFile(tokenFile, []byte(token), 0600)
	if err != nil {
		log.Printf("Failed to store token in %s: %v\n", tokenFile, err)
		return
	}
	log.Printf("Stored updated token %s in %s\n", token, tokenFile)
}
