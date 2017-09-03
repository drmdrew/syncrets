package vault

import (
	//	"errors"
	"fmt"
	"net/url"
	//	"io/ioutil"
	//	"log"
	//	"net/url"
	//	"strings"

	vaultapi "github.com/hashicorp/vault/api"
	//	"github.com/spf13/cobra"
	//	"github.com/spf13/viper"
)

// VaultReader is just the Read portion of the Vault client API
type VaultReader interface {
	Read(path string) (*vaultapi.Secret, error)
}

// ClientAPI is a composite API of all the Vault client APIs as interfaces
type ClientAPI interface {
	VaultReader
	Write(path string, data map[string]interface{}) (*vaultapi.Secret, error)
	SetToken(token string)
	Prompt(prompt string) string
}

type VaultClient struct {
	client *vaultapi.Client
}

func (vc *VaultClient) Read(path string) (*vaultapi.Secret, error) {
	return vc.client.Logical().Read(path)
}

func (vc *VaultClient) Write(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
	return vc.client.Logical().Write(path, data)
}

func (vc *VaultClient) SetToken(token string) {
	vc.client.SetToken(token)
}

func (vc *VaultClient) Prompt(prompt string) string {
	// TODO: find a better solution to prompt for token
	var token string
	fmt.Printf(prompt)
	fmt.Scanf("%s", &token)
	return token
}

// NewVaultClient creates a vaultClient using the supplied URL
func NewVaultClient(src *url.URL) (*VaultClient, error) {
	config := vaultapi.DefaultConfig()
	config.Address = fmt.Sprintf("%s://%s", "http", src.Host)
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	vc := &VaultClient{client}
	return vc, nil
}
