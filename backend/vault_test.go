package backend

import (
	"net/url"
	"testing"

	"github.com/spf13/viper"
)

func getViper(file string) *viper.Viper {
	v := viper.New()
	v.SetConfigFile(file)
	v.ReadInConfig()
	return v
}

func setupVault(t *testing.T) (*Vault, *mockVaultClient) {
	testViper := getViper("../testdata/syncrets-test1.yml")
	u, err := url.Parse("http://vault-a")
	if err != nil {
		t.Fatal(err)
	}
	mockVault := &mockVaultClient{}
	mockVault.data = make(map[string]map[string]interface{}, 1)
	endpoint := &Endpoint{}
	endpoint.Name = "vault-a"
	endpoint.RawURL = u
	endpoint.ServerURL = u
	v, err := NewVault(testViper, endpoint)
	if err != nil {
		t.Fatal(err)
	}
	v.client = mockVault
	return v, mockVault
}

func TestAuthenticate_withValidToken(t *testing.T) {
	v, mockVault := setupVault(t)
	mockVault.data["/auth/token/lookup-self"] = map[string]interface{}{"id": "mock-token"}
	err := v.Authenticate()
	if err != nil {
		t.Fatal(err)
	}
}
