package cmd

import (
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
)

var Debug bool

// RootCmd is the root cobra command for syncrets
var RootCmd = &cobra.Command{
	Use:   "subcommand [args] ...",
	Short: "subcommand required such as: auth, list, rm, sync",
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "debug logging output")
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	if !Debug {
		log.SetOutput(ioutil.Discard)
	}
}

// Execute the RootCmd
func Execute() {
	RootCmd.Execute()
}
