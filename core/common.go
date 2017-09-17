package core

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/drmdrew/syncrets/backend"
	"github.com/spf13/viper"
)

func parseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("Cannot parse URL: %v\n", s)
		return nil
	}
	log.Printf("Parsed URL: %v\n", u)
	return u
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

func resolveArgs(v *viper.Viper, args []string) (*backend.Endpoint, error) {
	if len(args) < 1 {
		return nil, errors.New("source argument is missing")
	}
	endpoint := &backend.Endpoint{}
	endpoint.RawURL = parseURL(args[0])
	if endpoint.RawURL == nil {
		return nil, errors.New("cannot parse url")
	}
	endpoint.Path = endpoint.RawURL.Path
	endpoint.Name = endpoint.RawURL.Hostname()
	u := resolveAlias(v, endpoint.Name)
	if u != nil {
		log.Printf("using alias: %v\n", u)
		endpoint.ServerURL = u
	} else {
		endpoint.ServerURL = endpoint.RawURL
	}
	log.Printf("using endpoint: %v\n", endpoint)
	return endpoint, nil
}

// NewVaultBackend returns a vault backend based on the supplied arguments
func NewVaultBackend(args []string) *backend.Endpoint {
	var endpoint *backend.Endpoint
	var vault *backend.Vault
	var err error
	if endpoint, err = resolveArgs(viper.GetViper(), args); err != nil {
		log.Fatal(err)
	}
	if vault, err = backend.NewVault(viper.GetViper(), endpoint); err != nil {
		log.Fatal(err)
	}
	if err = vault.Authenticate(); err != nil {
		log.Fatalf("Authenication failed: %v", err)
	}
	if !vault.IsValid() {
		log.Fatal("Authentication has failed!")
	}
	log.Print("Authentication was successful")
	vault.Store()
	endpoint.Vault = vault
	return endpoint
}
