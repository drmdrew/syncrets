package backend

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/drmdrew/syncrets/vault"
	"github.com/spf13/viper"
)

// Vault implements the syncrets vault backend
type Vault struct {
	hostname string
	url      *url.URL
	token    string
	viper    *viper.Viper
	client   vault.ClientAPI
}

// NewVault returns a vault instance
func NewVault(viper *viper.Viper, name string, url *url.URL) (*Vault, error) {
	v := &Vault{}
	v.viper = viper
	client, err := vault.NewClient(url)
	if err != nil {
		return nil, err
	}
	v.hostname = name
	v.client = client
	return v, nil
}

// GetClient returns a vault.ClientAPI
func (v *Vault) GetClient() vault.ClientAPI {
	return v.client
}

// prompt the user for information
func (v *Vault) prompt(prompt string) string {
	log.Printf("Prompting for user input: %s\n", prompt)
	var token string
	fmt.Printf(prompt)
	fmt.Scanf("%s", &token)
	return token
}

// Authenticate with the backend vault server
func (v *Vault) Authenticate() error {
	v.Load()
	if v.IsValid() {
		return nil
	}
	// re-authenticate if loaded token is invalid
	vkey := fmt.Sprintf("vault.%s.auth.method", v.hostname)
	method := v.viper.GetString(vkey)
	switch method {
	case "token":
		v.tokenAuth()
	case "userpass":
		v.userpassAuth()
	default:
		return fmt.Errorf("No valid auth.method configured for '%s'", v.hostname)
	}
	return nil
}

func (v *Vault) tokenAuth() {
	token := v.prompt("token: ")
	v.client.SetToken(token)
}

func (v *Vault) userpassAuth() {
	vkey := fmt.Sprintf("vault.%s.auth.username", v.hostname)
	username := v.viper.GetString(vkey)
	password := v.prompt("password: ")
	err := v.client.UserpassLogin(username, password)
	if err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}
}

// IsValid checks if the session with the backend is still valid
func (v *Vault) IsValid() bool {
	// use lookup-self to verify token is valid
	secret, err := v.client.Read("auth/token/lookup-self")
	if err != nil {
		log.Printf("lookup-self failed: %v\n", err)
		return false
	}
	id := secret.Data["id"]
	log.Printf("data: %v\n", secret.Data)
	log.Printf("token id: %v\n", id)
	return id != nil
}

func (v *Vault) Load() (string, error) {
	// load vault token from token.file if one is present
	vkey := fmt.Sprintf("vault.%s.token.file", v.hostname)
	tokenFile := v.viper.GetString(vkey)
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
	v.client.SetToken(token)
	return token, nil
}

func (v *Vault) Store() {
	// store the vault token in token.file if one is present
	vkey := fmt.Sprintf("vault.%s.token.file", v.hostname)
	tokenFile := v.viper.GetString(vkey)
	token := v.client.GetToken()
	err := ioutil.WriteFile(tokenFile, []byte(token), 0600)
	if err != nil {
		log.Printf("Failed to store token in %s: %v\n", tokenFile, err)
		return
	}
	log.Printf("Stored updated token %s in %s\n", token, tokenFile)
}
