package backend

import (
	"bytes"
	"strings"
	"testing"

	"github.com/drmdrew/syncrets/core"
)

var jsonTests = []struct {
	secrets  []core.Secret
	expected string
}{
	{[]core.Secret{core.Secret{"secret/citizen", "four"}},
		`{"secret":{"citizen":"four"}}`},
	{[]core.Secret{core.Secret{"secret/citizen/kane", "Rosebud"}},
		`{"secret":{"citizen":{"kane":"Rosebud"}}}`},
	{[]core.Secret{core.Secret{"secret/citizen", "four"}, core.Secret{"secret/citizen/kane", "Rosebud"}},
		`{"secret":{"citizen":{".":"four","kane":"Rosebud"}}}`},
}

func TestJSON_Marshal(t *testing.T) {
	for _, testcase := range jsonTests {
		buf := new(bytes.Buffer)
		j := NewJSONEndpoint()
		for _, s := range testcase.secrets {
			j.Visit(s)
		}
		j.Marshal(buf)
		result := strings.TrimSpace(buf.String())
		if result != testcase.expected {
			t.Fail()
			t.Fatalf("Expected: '%s' but result was: '%s'\n", testcase.expected, result)
		}
	}
}
