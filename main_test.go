package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestInitViper(t *testing.T) {
	v := initViper(viper.New(), "syncrets-test1", "testdata/")
	assert.NotNil(t, v.Sub("vault"))
	assert.Equal(t, v.Get("vault.vault-a.url"), "http://localhost:8201")
	assert.Equal(t, v.Get("vault.vault-a.auth.method"), "token")
	assert.Equal(t, v.Get("vault.vault-a.token.file"), "vault-a-token")
	assert.Equal(t, v.Get("vault.vault-b.url"), "http://localhost:8202")
}

func TestGetConfigFile(t *testing.T) {
	homeOrig := os.Getenv("HOME")
	os.Setenv("HOME", "testdata")
	assert.Equal(t, "testdata/.syncrets", getConfigFile())
	os.Setenv("HOME", homeOrig)
}
