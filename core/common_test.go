package core

import (
	"net/url"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func getViper(file string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(file)
	v.ReadInConfig()
	return v
}

var reverseLookupAliasTests = []struct {
	url    string
	expect string
}{
	{"http://localhost:8201", "vault-a"},
	{"https://localhost:8201", ""},   // wrong scheme
	{"https://unknown-tls", ""},      // unknown URL
	{"https://unknown-tls:8200", ""}, // unknown URL
}

func TestReverseLookupAlias_Valid(t *testing.T) {
	v := getViper("../testdata/syncrets-test1.yml")
	for _, tc := range reverseLookupAliasTests {
		u, _ := url.Parse(tc.url)
		alias := ReverseLookupAlias(v, u)
		assert.Equal(t, tc.expect, alias)
	}
}
