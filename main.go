package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/drmdrew/syncrets/cmd"
	"github.com/spf13/viper"
)

func main() {
	initViper(viper.GetViper(), "syncrets", ".")
	initCobra()
}

func initCobra() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func initViper(v *viper.Viper, name string, configDir string) *viper.Viper {
	v.SetConfigName(name)
	v.AddConfigPath(configDir)
	v.AddConfigPath(filepath.Join(os.Getenv("HOME"), ".syncrets"))
	err := v.ReadInConfig()
	if err != nil {
		log.Printf("No config file found %s. Using defaults\n", name)
	}
	return v
}
