package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/drmdrew/syncrets/backend"
	"github.com/drmdrew/syncrets/core"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	RootCmd.AddCommand(syncCmd)
}

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync secrets from vault",
	Long:  `Sync secrets from vault`,
	Run: func(cmd *cobra.Command, args []string) {
		srcArgs := args[0:1]
		src, err := backend.NewVaultBackend(viper.GetViper(), srcArgs)
		if err != nil {
			log.Fatal(err)
		}
		dstArgs := args[1:2]
		if strings.HasSuffix(dstArgs[0], ".json") {
			sync := backend.NewJSONEndpoint()
			src.Walk(sync)
			sync.Marshal(os.Stdout)
		} else {
			dst, err := backend.NewVaultBackend(viper.GetViper(), dstArgs)
			if err != nil {
				log.Fatal(err)
			}
			sync := &syncer{os.Stdout, dst}
			src.Walk(sync)
		}
	},
}

type syncer struct {
	out io.Writer
	dst core.Endpoint
}

func (sync *syncer) Visit(s core.Secret) {
	err := sync.dst.Write(s)
	fmt.Fprintf(sync.out, "%s => %s (%v)\n", s.Path, s.Path, err)
}
