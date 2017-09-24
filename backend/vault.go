package backend

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"github.com/drmdrew/syncrets/core"
	"github.com/spf13/viper"

	vaultapi "github.com/hashicorp/vault/api"
)

// Vault implements the syncrets vault backend
type Vault struct {
	name    string
	url     *url.URL
	origURL *url.URL
	path    string
	token   string
	viper   *viper.Viper
	client  VaultAPI
}

// SecretsReader is just the Read portion of the Vault client API
type SecretsReader interface {
	Read(path string) (*vaultapi.Secret, error)
}

// VaultAPI is a composite API of all the Vault client APIs as interfaces
type VaultAPI interface {
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

var newClientFunc = NewVaultClient

// NewClient creates a vaultClient using the supplied URL
func NewVaultClient(src *url.URL) (VaultAPI, error) {
	log.Printf("NewClient: making a *real* vault client...")
	config := vaultapi.DefaultConfig()
	config.Address = fmt.Sprintf("%s://%s", src.Scheme, src.Host)
	client, err := vaultapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	vc := &Client{client}
	return vc, nil
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

// GetName ...
func (v *Vault) GetName() string {
	return v.name
}

// GetPath ...
func (v *Vault) GetPath() string {
	return v.path
}

// GetClient returns a VaultAPI
func (v *Vault) GetClient() VaultAPI {
	return v.client
}

// GetClient returns a VaultAPI
func (v *Vault) GetURL() *url.URL {
	return v.url
}

// GetClient returns a VaultAPI
func (v *Vault) GetRawURL() *url.URL {
	return v.origURL
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
	vkey := fmt.Sprintf("vault.%s.auth.method", v.name)
	method := v.viper.GetString(vkey)
	switch method {
	case "token":
		v.tokenAuth()
	case "userpass":
		v.userpassAuth()
	default:
		return fmt.Errorf("No valid auth.method configured for '%s'", v.name)
	}
	return nil
}

func (v *Vault) tokenAuth() {
	token := v.prompt("token: ")
	v.client.SetToken(token)
}

func (v *Vault) userpassAuth() {
	vkey := fmt.Sprintf("vault.%s.auth.username", v.name)
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
	vkey := fmt.Sprintf("vault.%s.token.file", v.name)
	tokenFile := v.viper.GetString(vkey)
	tokenBytes, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		log.Printf("Error reading file %v: %v\n", tokenFile, err)
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
	vkey := fmt.Sprintf("vault.%s.token.file", v.name)
	tokenFile := v.viper.GetString(vkey)
	token := v.client.GetToken()
	err := ioutil.WriteFile(tokenFile, []byte(token), 0600)
	if err != nil {
		log.Printf("Failed to store token in %s: %v\n", tokenFile, err)
		return
	}
	log.Printf("Stored updated token %s in %s\n", token, tokenFile)
}

func (src *Vault) Write(secret core.Secret) error {
	data := map[string]interface{}{
		"value": secret.Value,
	}
	_, err := src.GetClient().Write(secret.Path, data)
	return err
}

// Walk the secrets...
func (src *Vault) Walk(visitor core.Visitor) {
	var prefixes []string
	path := src.GetPath()
	prefixes = append(prefixes, path)
	for len(prefixes) > 0 {
		// pop a prefix from the front of the slice
		var prefix string
		prefix, prefixes = prefixes[0], prefixes[1:]
		secret, err := src.GetClient().List(prefix)
		if err != nil {
			continue
		}
		if secret != nil {
			for _, val := range secret.Data["keys"].([]interface{}) {
				s := val.(string)
				if strings.HasSuffix(s, "/") {
					// push a new prefix at the end of the slice
					prefixes = append(prefixes, prefix+s)
				} else {
					// this leaf has a secret value...
					// ... now print it
					sep := "/"
					if strings.HasSuffix(prefix, "/") {
						sep = ""
					}
					path := fmt.Sprintf("%s%s%s", prefix, sep, s)
					//fmt.Printf("%s\n", path)
					// ... so copy it to dst vault
					value, err := src.GetClient().Read(path)
					if value != nil {
						//fmt.Printf("   -> value.Data['value']: %s\n", value.Data["value"])
						data := value.Data
						//bValue, bErr := dst.Vault.GetClient().Write(path, data)
						secret := core.Secret{path, data["value"].(string)}
						visitor.Visit(secret)
						//fmt.Printf("   -> path=%s, secret=%v, err=%v\n", path, value, err)
					} else {
						fmt.Printf("   !! err: %v\n", err)
					}
				}
			}
		}
	}
}

func (v *Vault) resolveArgs(args []string) error {
	if len(args) < 1 {
		return errors.New("source argument is missing")
	}
	v.origURL = core.ParseURL(args[0])
	if v.origURL == nil {
		return errors.New("cannot parse url")
	}
	v.path = v.origURL.Path
	v.name = v.origURL.Hostname()
	u := core.ResolveAlias(v.viper, v.name)
	if u != nil {
		log.Printf("using alias: %v\n", u)
		v.url = u
	} else {
		v.url = v.origURL
	}
	log.Printf("%s using url: %v\n", v.name, v.url)
	return nil
}

// NewVaultBackend returns a vault backend based on the supplied arguments
func NewVaultBackend(viper *viper.Viper, args []string) (*Vault, error) {
	v := &Vault{}
	v.viper = viper
	if err := v.resolveArgs(args); err != nil {
		return nil, err
	}
	client, err := newClientFunc(v.url)
	if err != nil {
		return nil, err
	}
	v.client = client
	if err := v.Authenticate(); err != nil {
		log.Fatalf("Authenication failed: %v", err)
	}
	if !v.IsValid() {
		log.Fatal("Authentication has failed!")
	}
	log.Print("Authentication was successful")
	v.Store()
	return v, nil
}
