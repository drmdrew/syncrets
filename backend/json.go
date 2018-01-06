package backend

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/drmdrew/syncrets/core"
)

// JSONEndpoint ...
type JSONEndpoint struct {
	kv map[string]interface{}
}

// NewJSONEndpoint ...
func NewJSONEndpoint() *JSONEndpoint {
	return &JSONEndpoint{make(map[string]interface{})}
}

// Visit ...
func (j *JSONEndpoint) Visit(s core.Secret) {
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
func (j *JSONEndpoint) Marshal(out io.Writer) error {
	b, err := json.Marshal(j.kv)
	if err != nil {
		log.Print(err)
		return err
	}
	fmt.Fprintf(out, "%s\n", string(b[:]))
	return nil
}
