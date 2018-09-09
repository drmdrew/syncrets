// +build integration

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func execCommand(t *testing.T, cmdName string, cmdArgs []string) []string {
	var cmdOut []byte
	var err error
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintf(os.Stderr, "%v failed: %v\n%v", cmdName, err, string(cmdOut))
		t.FailNow()
	}
	scanner := bufio.NewScanner(bytes.NewReader(cmdOut))
	var output []string
	for scanner.Scan() {
		output = append(output, scanner.Text())
	}
	return output
}

/* libcompose doesn't work with docker-compose version 3 */
/*
func dockerComposeSetup() {
	project, err := docker.NewProject(&ctx.Context{
		Context: project.Context{
			ComposeFiles: []string{"docker-compose.yml"},
			ProjectName:  "syncrets",
		},
	}, nil)

	if err != nil {
		log.Fatal(err)
	}

	err = project.Up(context.Background(), options.Up{})

	if err != nil {
		log.Fatal(err)
	}
}
*/
/*
func dockerComposeSetup(t *testing.T) {
	output := execCommand(t, "docker-compose", []string{"up", "-d"})
	log.Println(output)
}
*/

func TestIntegration_SyncretsList(t *testing.T) {
	// dockerComposeSetup(t)
	expected := []string{
		"/secret/foo", "/secret/gilbert", "/secret/foo/bar", "/secret/it/was/the/best/of/times",
	}
	actual := execCommand(t, "./syncrets", []string{"list", "vault://vault-a/secret/"})
	for i, key := range actual {
		if key != expected[i] {
			fmt.Fprintf(os.Stderr, "%d: expected '%v', got '%v'", i, expected[i], key)
			t.FailNow()
		}
	}
}

func TestIntegration_SyncretsSyncToJSON(t *testing.T) {
	actual := execCommand(t, "./syncrets", []string{"sync", "vault://vault-a/secret/", "testoutput/testOutputSyncToJSON.json"})
	if len(actual) != 0 {
		t.Fatalf("Expected no stdout, got: %v", actual)
	}
	bytes, err := ioutil.ReadFile("testoutput/testOutputSyncToJSON.json")
	if err != nil {
		t.Fatalf("Error reading output file: %v", err)
	}
	json := string(bytes)
	jsonTrimmed := strings.TrimSpace(json)
	expected := `{"secret":{"foo":{".":"bar","bar":"foobar"},"gilbert":"sullivan","it":{"was":{"the":{"best":{"of":{"times":"it was the worst of times"}}}}}}}`
	if jsonTrimmed != expected {
		t.Fatalf("JSON output not as expected: %v", json)
	}
}
