package core

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/drmdrew/syncrets/backend"
	"github.com/spf13/viper"
)

// ...
func ParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("Cannot parse URL: %v\n", s)
		return nil
	}
	return u
}

// ...
func ResolveAlias(v *viper.Viper, alias string) *url.URL {
	vkey := fmt.Sprintf("vault.%s.url", alias)
	vurl := v.GetString(vkey)
	log.Printf("Checking for alias: %v", vkey)
	if vurl != "" {
		log.Printf("using alias: %s\n", vurl)
		return ParseURL(vurl)
	}
	return nil
}

// ...
func ResolveArgs(v *viper.Viper, args []string) (string, *url.URL, error) {
	if len(args) < 1 {
		return "", nil, errors.New("source argument is missing")
	}
	u := ParseURL(args[0])
	if u == nil {
		return "", nil, errors.New("cannot parse url")
	}
	hostname := u.Hostname()
	u = ResolveAlias(v, hostname)
	if u != nil {
		log.Printf("using alias: %v\n", u)
	}
	return hostname, u, nil
}

// NewVaultBackend ...
func NewVaultBackend(args []string) *backend.Vault {
	var hostname string
	var url *url.URL
	var vault *backend.Vault
	var err error
	if hostname, url, err = ResolveArgs(viper.GetViper(), args); err != nil {
		log.Fatal(err)
	}
	if vault, err = backend.NewVault(viper.GetViper(), hostname, url); err != nil {
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
	return vault
}
