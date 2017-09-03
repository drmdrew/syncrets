package vault

import (
	"fmt"
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
	Write(path string, data map[string]interface{}) (*vaultapi.Secret, error)
	SetToken(token string)
	Prompt(prompt string) string
}

// Client for communicating with vault backends
type Client struct {
	client *vaultapi.Client
}

// Read a secret from a vault backend
func (vc *Client) Read(path string) (*vaultapi.Secret, error) {
	return vc.client.Logical().Read(path)
}

// Write secrets to a vault backend
func (vc *Client) Write(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
	return vc.client.Logical().Write(path, data)
}

// Sets the token for authentication with a vault backend
func (vc *Client) SetToken(token string) {
	vc.client.SetToken(token)
}

// Prompt the user for information
func (vc *Client) Prompt(prompt string) string {
	// TODO: find a better solution to prompt for token
	var token string
	fmt.Printf(prompt)
	fmt.Scanf("%s", &token)
	return token
}

// NewClient creates a vaultClient using the supplied URL
func NewClient(src *url.URL) (*Client, error) {
	config := vaultapi.DefaultConfig()
	config.Address = fmt.Sprintf("%s://%s", "http", src.Host)
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	vc := &Client{client}
	return vc, nil
}
