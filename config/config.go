package config

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path"
)

var Config = config{
	MaxAlienMoves:  100000,
	MapFile:        "map.txt",
	AlienNamesFile: "alien-names.txt",
}

type config struct {
	MaxAlienMoves  uint
	MapFile        string
	AlienNamesFile string
	DebugMode      bool
}

func Path() string {
	filePath := "config.json"
	cwd, err := os.Getwd()
	if err == nil {
		return path.Join(cwd, filePath)
	}
	return filePath
}

func Load(file io.Reader) error {
	fileData, _ := ioutil.ReadAll(file)
	err := json.Unmarshal(fileData, &Config)
	if err != nil {
		return err
	}
	return nil
}

func Write(file io.Writer) error {
	configJson, err := json.MarshalIndent(Config, "", "    ")
	if err != nil {
		return err
	}
	_, err = file.Write(configJson)
	return err
}
