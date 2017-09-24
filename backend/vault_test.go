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

func setupVault(t *testing.T, mockData map[string]map[string]interface{}) (*Vault, *mockVaultClient) {
	testViper := getViper("./testdata/syncrets-test1.yml")
	mockVault := &mockVaultClient{}
	mockVault.data = mockData
	newClientFunc = func(src *url.URL) (VaultAPI, error) {
		return mockVault, nil
	}
	args := []string{"http://vault-a"}
	v, err := NewVaultBackend(testViper, args)
	if err != nil {
		t.Fatal(err)
	}
	return v, mockVault
}

func TestAuthenticate_withValidToken(t *testing.T) {
	mockData := make(map[string]map[string]interface{}, 1)
	mockData["auth/token/lookup-self"] = map[string]interface{}{"id": "mock-token"}
	v, mockVault := setupVault(t, mockData)
	t.Logf("test configured mock vault with data: %v\n", mockVault.data)
	err := v.Authenticate()
	if err != nil {
		t.Fatal(err)
	}
}
