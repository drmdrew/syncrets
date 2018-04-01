package cmd

import (
	"bufio"
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

func createFileAndWriter(s string) (*os.File, *bufio.Writer) {
	f, err := os.Create(s)
	if err != nil {
		log.Fatal(err)
	}
	w := bufio.NewWriter(f)
	return f, w
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
			f, w := createFileAndWriter(dstArgs[0])
			defer f.Close()
			sync := backend.NewJSONEndpoint()
			src.Walk(sync)
			sync.Marshal(w)
			w.Flush()
		} else if strings.HasSuffix(dstArgs[0], ".ejson") {
			f, w := createFileAndWriter(dstArgs[0])
			defer f.Close()
			sync := backend.NewEJSONEndpoint()
			src.Walk(sync)
			sync.Marshal(w)
			w.Flush()
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
