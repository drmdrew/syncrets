package cmd

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestAuth_LocalhostURL(t *testing.T) {
	args := []string{"http://localhost:8200"}

	auth, _ := newAuthenticator(viper.GetViper(), authCmd, args)
	assert.Equal(t, "localhost:8200", auth.url.Host)
}
func TestAuth_ResolvesAlias(t *testing.T) {
	args := []string{"http://vault-a:8200"}
	testViper := viper.New()
	testViper.SetConfigFile("../testdata/syncrets-test1.yml")
	testViper.ReadInConfig()

	auth, _ := newAuthenticator(testViper, authCmd, args)
	assert.Equal(t, "localhost:8201", auth.url.Host)
}

type fakeVaultReader struct {
	data map[string]interface{}
}

func (fvr *fakeVaultReader) Read(path string) (*vaultapi.Secret, error) {
	s := &vaultapi.Secret{}
	s.Data = fvr.data //make(map[string]interface{})
	return s, nil
}

func TestAuth_IsValid(t *testing.T) {
	args := []string{"http://localhost:8200"}
	auth, err := newAuthenticator(viper.GetViper(), authCmd, args)
	fvr := &fakeVaultReader{make(map[string]interface{})}
	auth.reader = fvr
	assert.Nil(t, err)
	assert.True(t, auth.isValid())

}
