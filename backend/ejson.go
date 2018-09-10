package backend

import (
	"bytes"
	"encoding/json"
	"io"
	"log"

	"github.com/Shopify/ejson"
	"github.com/drmdrew/syncrets/core"
	"github.com/spf13/viper"
)

type EJSONEndpoint struct {
	kv map[string]interface{}
}

// NewEJSONEndpoint ...
func NewEJSONEndpoint() *EJSONEndpoint {
	return &EJSONEndpoint{make(map[string]interface{})}
}

// Visit ...
func (j *EJSONEndpoint) Visit(s core.Secret) {
	AddSecretToKV(s, j.kv)
}

// Marshal ...
func (j *EJSONEndpoint) Marshal(out io.Writer) error {
	var jsonBytes []byte
	var err error
	ejsonPubkey := viper.GetString("ejson.public_key")
	j.kv["_public_key"] = ejsonPubkey
	if jsonBytes, err = json.Marshal(j.kv); err != nil {
		log.Printf("ERROR: %v\n", err)
		return err
	}
	in := bytes.NewReader(jsonBytes)
	_, err = ejson.Encrypt(in, out)
	if err != nil {
		log.Printf("ERROR: %v\n", err)
		return err
	}
	return nil
}
