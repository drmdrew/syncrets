package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/drmdrew/syncrets/cmd"
	"github.com/spf13/viper"
)

func main() {
	initViper(viper.GetViper(), "syncrets", getConfigFile())
	initCobra()
}

func initCobra() {
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func getConfigFile() string {
	return filepath.Join(os.Getenv("HOME"), ".syncrets")
}

func initViper(v *viper.Viper, name string, path string) *viper.Viper {
	v.SetConfigName(name)
	v.AddConfigPath(path)
	err := v.ReadInConfig()
	if err != nil {
		log.Printf("No config file found %s on path %s. Using defaults\n", name, path)
	}
	return v
}
