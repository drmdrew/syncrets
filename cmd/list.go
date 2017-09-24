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
	RootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List secrets from vault",
	Long:  `List secrets from vault`,
	Run: func(cmd *cobra.Command, args []string) {
		list := &lister{os.Stdout}
		srcArgs := args[0:1]
		src, err := backend.NewVaultBackend(viper.GetViper(), srcArgs)
		if err != nil {
			log.Fatal(err)
		}
		src.Walk(list)
	},
}

type lister struct {
	out io.Writer
}

func (l *lister) Visit(s core.Secret) {
	fmt.Fprintf(l.out, "%s\n", s.Path)
}
