package backend

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strings"
	"syscall"

	"github.com/drmdrew/syncrets/core"
	"github.com/spf13/viper"

	vaultapi "github.com/hashicorp/vault/api"
	"golang.org/x/crypto/ssh/terminal"
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
	isValid *bool
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
	Delete(path string) (*vaultapi.Secret, error)
	GetToken() string
	SetToken(token string)
}

// Client for communicating with vault backends
type Client struct {
	client *vaultapi.Client
}

var newClientFunc = NewVaultClient

// NewVaultClient creates a vaultClient using the supplied URL
func NewVaultClient(src *url.URL) (VaultAPI, error) {
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

// Delete secret from a vault backend
func (vc *Client) Delete(path string) (*vaultapi.Secret, error) {
	return vc.client.Logical().Delete(path)
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

// GetURL ...
func (v *Vault) GetURL() *url.URL {
	return v.url
}

// GetRawURL ...
func (v *Vault) GetRawURL() *url.URL {
	return v.origURL
}

// prompt the user for information
func (v *Vault) prompt(prompt string) (string, error) {
	log.Printf("Prompting for user input: %s\n", prompt)
	fmt.Printf(prompt)
	bytes, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Printf("Error reading response to prompt: %v\n", err)
		return "", err
	}
	return string(bytes), nil
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
		log.Printf("Unknown auth.method '%s' configured for '%s'\n", method, v.name)
		if _, hasEnv := os.LookupEnv("VAULT_TOKEN"); hasEnv {
			log.Printf("Using VAULT_TOKEN environment variable\n")
			v.envAuth()
		} else {
			log.Printf("Defaulting to 'token' authentication\n")
			v.tokenAuth()
		}
	}
	return nil
}

func (v *Vault) envAuth() {
	token := os.Getenv("VAULT_TOKEN")
	v.client.SetToken(token)
}

func (v *Vault) tokenAuth() error {
	token, err := v.prompt("token: ")
	if err == nil {
		v.client.SetToken(token)
	}
	return nil
}

func (v *Vault) userpassAuth() error {
	vkey := fmt.Sprintf("vault.%s.auth.username", v.name)
	username := v.viper.GetString(vkey)
	password, err := v.prompt("password: ")
	if err == nil {
		err = v.client.UserpassLogin(username, password)
	}
	if err != nil {
		log.Printf("Authentication failed: %v", err)
		return err
	}
	return nil
}

// IsValid checks if the session with the backend is still valid
func (v *Vault) IsValid() bool {
	if v.isValid != nil {
		return *v.isValid
	}
	// use lookup-self to verify token is valid
	valid := false
	secret, err := v.client.Read("auth/token/lookup-self")
	if err != nil {
		log.Printf("lookup-self failed: %v\n", err)
	} else {
		id := secret.Data["id"]
		valid = id != nil
		log.Printf("lookup-self returned %t, accessor: %v\n", valid, secret.Data["accessor"])
	}
	v.isValid = &valid
	return *v.isValid
}

// Load ...
func (v *Vault) Load() (string, error) {
	// load vault token from token.file if one is present
	vkey := fmt.Sprintf("vault.%s.token.file", v.name)
	tokenFile := v.viper.GetString(vkey)
	if tokenFile == "" {
		return "", fmt.Errorf("No token defined for %v", vkey)
	}
	tokenBytes, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		log.Printf("Error reading file %v: %v\n", tokenFile, err)
		return "", err
	}
	log.Printf("%v is configured: %v\n", vkey, tokenFile)
	var token string
	if tokenBytes != nil {
		token = strings.TrimSpace(string(tokenBytes))
	} else {
		return "", fmt.Errorf("Unable to read token from %s", tokenFile)
	}
	v.client.SetToken(token)
	return token, nil
}

// Store ...
func (v *Vault) Store() {
	// store the vault token in token.file if one is present
	vkey := fmt.Sprintf("vault.%s.token.file", v.name)
	tokenFile := v.viper.GetString(vkey)
	if tokenFile == "" {
		log.Printf("Not storing token. No token file configured for %s\n", vkey)
		return
	}
	token := v.client.GetToken()
	err := ioutil.WriteFile(tokenFile, []byte(token), 0600)
	if err != nil {
		log.Printf("Failed to write token file '%s': %v\n", tokenFile, err)
		return
	}
	log.Printf("Stored updated token in %s\n", tokenFile)
}

// Write ...
func (src *Vault) Write(secret core.Secret) error {
	data := map[string]interface{}{
		"value": secret.Value,
	}
	_, err := src.GetClient().Write(secret.Path, data)
	return err
}

// Delete ...
func (src *Vault) Delete(secret core.Secret) error {
	_, err := src.GetClient().Delete(secret.Path)
	return err
}

// Walk the secrets...
func (src *Vault) Walk(visitor core.Visitor) {
	var prefixes []string
	path := src.GetPath()
	prefixes = append(prefixes, path)
	log.Printf("-> walk prefixes: %v\n", prefixes)
	for len(prefixes) > 0 {
		// pop a prefix from the front of the slice
		var prefix string
		prefix, prefixes = prefixes[0], prefixes[1:]
		secret, err := src.GetClient().List(prefix)
		if err != nil {
			log.Printf("   -> list error: %v\n", err)
			continue
		}
		log.Printf("   -> list prefix %v: %v\n", prefix, secret != nil)
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
						data := value.Data
						//bValue, bErr := dst.Vault.GetClient().Write(path, data)
						secret := core.Secret{path, data["value"].(string)}
						visitor.Visit(secret)
						log.Printf("       <- visited path=%s, err=%v\n", path, err)
					} else {
						fmt.Printf("       !! err: %v\n", err)
					}
				}
			}
		}
		// check if prefix itself is a leaf and has a secret value
		leafSecret, leafErr := src.readSecret(prefix)
		if leafErr != nil {
			log.Printf("   -> readSecret prefix %v error: %v\n", prefix, leafErr)
			continue
		}
		if leafSecret != nil {
			log.Printf("   -> VISIT leaf prefix %v -> %v\n", prefix, leafSecret.Path)
			visitor.Visit(*leafSecret)
		}
	}
}

func (v *Vault) readSecret(path string) (*core.Secret, error) {
	value, err := v.GetClient().Read(path)
	var secret *core.Secret
	if err == nil && value != nil {
		data := value.Data
		secret = &core.Secret{path, data["value"].(string)}
	}
	return secret, err
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
	alias := core.ReverseLookupAlias(v.viper, v.origURL)
	if alias == "" {
		alias = v.origURL.Hostname()
	}
	u := core.ResolveAlias(v.viper, alias)
	if u != nil {
		v.url = u
	} else {
		v.url = v.origURL
	}
	v.name = alias
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
		log.Printf("Authenication failed: %v", err)
		return nil, err
	}
	if !v.IsValid() {
		log.Print("Authentication has failed!")
		return nil, err
	}
	log.Print("Authentication was successful")
	v.Store()
	return v, nil
}
