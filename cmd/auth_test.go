package cmd

import (
	"testing"

	vaultapi "github.com/hashicorp/vault/api"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type fakeVaultClient struct {
	data map[string]interface{}
}

func (fvr *fakeVaultClient) Read(path string) (*vaultapi.Secret, error) {
	s := &vaultapi.Secret{}
	s.Data = fvr.data //make(map[string]interface{})
	return s, nil
}

func (fvr *fakeVaultClient) List(path string) (*vaultapi.Secret, error) {
	s := &vaultapi.Secret{}
	s.Data = fvr.data //make(map[string]interface{})
	return s, nil
}

func (fvr *fakeVaultClient) Write(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
	s := &vaultapi.Secret{}
	s.Data = fvr.data //make(map[string]interface{})
	return s, nil
}

func (fvr *fakeVaultClient) SetToken(token string) {
	// do nothing
}

func (fvr *fakeVaultClient) GetToken() string {
	return ""
}

func (fvr *fakeVaultClient) UserpassLogin(username string, password string) error {
	return nil
}

func (fvr *fakeVaultClient) TokenIsValid() bool {
	return true
}

func getViper(file string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(file)
	v.ReadInConfig()
	return v
}

func TestAuth_LocalhostURL(t *testing.T) {
	args := []string{"http://localhost:8200"}
	auth, _ := newAuthenticator(viper.GetViper(), authCmd, args)
	assert.Equal(t, "localhost:8200", auth.url.Host)
}
func TestAuth_ResolvesAlias(t *testing.T) {
	args := []string{"http://vault-a:8200"}
	testViper := getViper("../testdata/syncrets-test1.yml")
	auth, _ := newAuthenticator(testViper, authCmd, args)
	assert.Equal(t, "localhost:8201", auth.url.Host)
}

func TestAuth_IsValid(t *testing.T) {
	args := []string{"http://localhost:8200"}
	auth, err := newAuthenticator(viper.GetViper(), authCmd, args)
	fvc := &fakeVaultClient{make(map[string]interface{})}
	auth.client = fvc
	assert.Nil(t, err)
	assert.True(t, auth.isValid())
}
