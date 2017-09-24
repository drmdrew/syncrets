package core

import (
	"fmt"
	"log"
	"net/url"

	"github.com/spf13/viper"
)

// ParseURL to parse a URL
func ParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("Cannot parse URL: %v\n", s)
		return nil
	}
	log.Printf("Parsed URL: %v\n", u)
	return u
}

// ResolveAlias to resolve an alias
func ResolveAlias(v *viper.Viper, alias string) *url.URL {
	vkey := fmt.Sprintf("vault.%s.url", alias)
	vurl := v.GetString(vkey)
	log.Printf("Checking for alias: %v", vkey)
	if vurl != "" {
		log.Printf("Found an alias: %s\n", vurl)
		return ParseURL(vurl)
	}
	return nil
}
