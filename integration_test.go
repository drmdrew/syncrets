// +build integration
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestIntegration_SyncretsList(t *testing.T) {
	expected := []string{
		"/secret/foo", "/secret/gilbert", "/secret/foo/bar", "/secret/it/was/the/best/of/times"}
	cmdName := "./syncrets"
	cmdArgs := []string{"list", "vault://vault-a/secret/"}
	var cmdOut []byte
	var err error
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintf(os.Stderr, "There was an error running syncrets: %v\n%v", err, string(cmdOut))
		t.FailNow()
	}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOut))
	var actual []string
	for scanner.Scan() {
		actual = append(actual, scanner.Text())
	}
	for i, key := range actual {
		if key != expected[i] {
			fmt.Fprintf(os.Stderr, "line %d: expected '%v', got '%v'", expected[i], key)
			t.FailNow()
		}
	}
}
