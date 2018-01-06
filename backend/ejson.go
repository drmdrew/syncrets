package backend

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/Shopify/ejson"
	"github.com/drmdrew/syncrets/core"
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
	steps := strings.Split(s.Path, "/")
	var kv = j.kv
	for _, step := range steps[:len(steps)-1] {
		if step == "" {
			continue
		}
		if _, ok := kv[step]; !ok {
			kv[step] = make(map[string]interface{})
		}
		if m, ok := kv[step].(map[string]interface{}); !ok {
			m = make(map[string]interface{})
			m["."] = kv[step]
			kv[step] = m
			kv = m
		} else {
			kv = m
		}
	}
	lastStep := steps[len(steps)-1]
	kv[lastStep] = s.Value
}

// Marshal ...
func (j *EJSONEndpoint) Marshal(out io.Writer) error {
	var jsonBytes []byte
	var err error
	// TODO: configurable ejson public key
	j.kv["_public_key"] = "a9d52487a1232e5c292a9680f4a44a84ea302ba05ff12d2e9d11662d20fc0139"
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
