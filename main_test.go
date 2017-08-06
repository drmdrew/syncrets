package main

import (
	"io/ioutil"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestInitViper(t *testing.T) {
	v := initViper("syncrets-test1", "testdata/")
	assert.NotNil(t, v.Sub("vault"))
	assert.Equal(t, v.Get("vault.vault-a.url"), "http://localhost:8201")
	assert.Equal(t, v.Get("vault.vault-b.url"), "http://localhost:8202")
}
