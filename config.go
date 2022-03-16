package aliens

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"io/ioutil"
	"os"
	"path"
)

const ConfigFile = "config.json"

var fs = afero.NewOsFs()

var Config = config{
	MaxAlienMoves:       100000,
	MapFile:             "map.txt",
	AlienNamesFile:      "alien-names.txt",
}

type config struct {
	MaxAlienMoves       uint
	MapFile             string
	AlienNamesFile      string
	DebugMode           bool
}

func getConfigFilePath() string {
	cwd, err := os.Getwd()
	if err == nil {
		return path.Join(cwd, ConfigFile)
	}
	return ConfigFile
}

func LoadConfig() error {
	configFile := getConfigFilePath()
	log.Printf("Loading config from: %s\n", configFile)
	file, err := fs.Open(configFile)
	if err != nil {
		return err
	}
	fileData, _ := ioutil.ReadAll(file)
	err = json.Unmarshal(fileData, &Config)
	if err != nil {
		return err
	}
	return nil
}

func WriteConfig() error {
	configFile := getConfigFilePath()
	log.Infof("Writing config to: %s\n", configFile)
	configJson, err := json.MarshalIndent(Config, "", "    ")
	if err != nil {
		return err
	}
	return afero.WriteFile(fs, configFile, configJson, 0644)
}
