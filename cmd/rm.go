package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/drmdrew/syncrets/backend"
	"github.com/drmdrew/syncrets/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(rmCmd)
}

var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove secrets from vault",
	Long:  `Remove secrets from vault`,
	Run: func(cmd *cobra.Command, args []string) {
		srcArgs := args[0:1]
		src, err := backend.NewVaultBackend(viper.GetViper(), srcArgs)
		if err != nil {
			log.Fatal(err)
		}
		rm := &remover{os.Stdout, src}
		src.Walk(rm)
	},
}

type remover struct {
	out      io.Writer
	endpoint core.Endpoint
}

func (rm *remover) Visit(s core.Secret) {
	fmt.Fprintf(rm.out, "Deleted %s\n", s.Path)
	rm.endpoint.Delete(s)
}
