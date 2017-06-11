package main

import (
	"log"

	"github.com/drmdrew/syncrets/cmd"
)

func main() {

	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}
