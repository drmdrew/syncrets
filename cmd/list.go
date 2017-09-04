package cmd

import (
	"fmt"
	"log"

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
		auth, err := newAuthenticator(viper.GetViper(), cmd, args)
		if err != nil {
			log.Fatal(err)
		}
		if err := auth.authenticate(); err != nil {
			log.Fatalf("Authenication failed: %v", err)
		}
		if !auth.isValid() {
			log.Fatal("Authentication has failed!")
		}
		log.Print("Authentication was successful")
		auth.store()
		secret, err := auth.client.List("secret/")
		for k, v := range secret.Data {
			fmt.Printf("%v, %v\n", k, v)
		}
	},
}
