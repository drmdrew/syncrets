package vault

import (
	"fmt"
	"log"
	"net/url"

	vaultapi "github.com/hashicorp/vault/api"
)

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
