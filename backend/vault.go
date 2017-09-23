package backend

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/spf13/viper"

	vaultapi "github.com/hashicorp/vault/api"
)

// Vault implements the syncrets vault backend
type Vault struct {
	hostname string
	url      *url.URL
	token    string
	viper    *viper.Viper
	client   ClientAPI
}

// SecretsReader is just the Read portion of the Vault client API
type SecretsReader interface {
	Read(path string) (*vaultapi.Secret, error)
}

// ClientAPI is a composite API of all the Vault client APIs as interfaces
type ClientAPI interface {
	SecretsReader
	List(path string) (*vaultapi.Secret, error)
	UserpassLogin(username string, password string) error
	TokenIsValid() bool
	Write(path string, data map[string]interface{}) (*vaultapi.Secret, error)
	GetToken() string
	SetToken(token string)
}

// Client for communicating with vault backends
type Client struct {
	client *vaultapi.Client
}

// Read a secret from a vault backend
func (vc *Client) Read(path string) (*vaultapi.Secret, error) {
	return vc.client.Logical().Read(path)
}

// List secrets from a vault backend
func (vc *Client) List(path string) (*vaultapi.Secret, error) {
	secret, err := vc.client.Logical().List(path)
	//	if secret != nil {
	//		log.Printf("secret.Data: %v\n", secret.Data)
	//	}
	return secret, err
}

// Write secrets to a vault backend
func (vc *Client) Write(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
	return vc.client.Logical().Write(path, data)
}

// SetToken sets the token for authentication with a vault backend
func (vc *Client) SetToken(token string) {
	vc.client.SetToken(token)
}

// GetToken returns the current authentication token being used
func (vc *Client) GetToken() string {
	return vc.client.Token()
}

// UserpassLogin performs a username+password login with vault
func (vc *Client) UserpassLogin(username string, password string) error {
	data := map[string]interface{}{
		"password": password,
	}
	url := fmt.Sprintf("auth/userpass/login/%s", username)
	result, err := vc.Write(url, data)
	if err != nil {
		return err
	}
	vc.SetToken(result.Auth.ClientToken)
	return nil
}

// TokenIsValid queries the vault server to verify that the token is valid
func (vc *Client) TokenIsValid() bool {
	// use lookup-self to verify token is valid
	_, err := vc.Read("auth/token/lookup-self")
	if err != nil {
		log.Printf("lookup-self failed: %v\n", err)
		return false
	}
	return true
}

// NewClient creates a vaultClient using the supplied URL
func NewClient(src *url.URL) (*Client, error) {
	config := vaultapi.DefaultConfig()
	config.Address = fmt.Sprintf("%s://%s", src.Scheme, src.Host)
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	vc := &Client{client}
	return vc, nil
}

// NewVault returns a vault instance
func NewVault(viper *viper.Viper, endpoint *Endpoint) (*Vault, error) {
	v := &Vault{}
	v.viper = viper
	client, err := NewClient(endpoint.ServerURL)
	if err != nil {
		return nil, err
	}
	v.hostname = endpoint.Name
	v.client = client
	return v, nil
}

// GetClient returns a ClientAPI
func (v *Vault) GetClient() ClientAPI {
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
