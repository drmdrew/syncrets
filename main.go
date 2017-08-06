package main

import (
	"log"

	"github.com/drmdrew/syncrets/cmd"
	"github.com/spf13/viper"
)

func main() {
	initCobra()
	initViper("syncrets", "/Users/drew.syncrets")
}

func initCobra() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initViper(name string, path string) *viper.Viper {
	v := viper.New()
	v.SetConfigName(name)
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		log.Printf("No config file found %s on path %s. Using defaults\n", name, path)
	}
	return v
}
