package backend

import (
	"log"

	vaultapi "github.com/hashicorp/vault/api"
)

type mockVaultClient struct {
	valid bool
	token string
	data  map[string]map[string]interface{}
}

func (v *mockVaultClient) Read(path string) (*vaultapi.Secret, error) {
	s := &vaultapi.Secret{}
	s.Data = v.data[path] //make(map[string]interface{})
	log.Printf("mock vault data: %v", v.data)
	log.Printf("mock vault client Read(%s) => %v", path, s.Data)
	return s, nil
}

func (v *mockVaultClient) List(path string) (*vaultapi.Secret, error) {
	s := &vaultapi.Secret{}
	s.Data = v.data[path] //make(map[string]interface{})
	return s, nil
}

func (v *mockVaultClient) Write(path string, data map[string]interface{}) (*vaultapi.Secret, error) {
	s := &vaultapi.Secret{}
	s.Data = v.data[path] //make(map[string]interface{})
	return s, nil
}

func (v *mockVaultClient) Delete(path string) (*vaultapi.Secret, error) {
	s := &vaultapi.Secret{}
	return s, nil
}

func (v *mockVaultClient) SetToken(token string) {
	v.token = token
}

func (v *mockVaultClient) GetToken() string {
	return v.token
}

func (v *mockVaultClient) UserpassLogin(username string, password string) error {
	return nil
}

func (v *mockVaultClient) TokenIsValid() bool {
	return v.valid
}
