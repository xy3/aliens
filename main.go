package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/xy3/aliens/cmd"
	"github.com/xy3/aliens/config"
	"os"
)

func main() {
	logger := log.New()
	configFile := config.Path()
	file, err := os.Open(configFile)
	if err != nil {
		log.WithError(err).Warnf("Failed to open config file [%s] for reading", configFile)
	} else {
		log.Printf("Loading config from: %s\n", configFile)
		err = config.Load(file)
		if err != nil {
			log.Println("Failed to find or read a 'config.json' file, using default config values or flag overrides")
		}
	}

	rootCmd := cmd.NewRootCmd(logger)
	err = rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
