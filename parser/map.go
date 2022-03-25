package parser

import (
	"github.com/xy3/aliens/world"
	"io"
	"io/ioutil"
	"strings"
)

// ParseMap reads in the map text file and parses the cities into the worldMap
func ParseMap(file io.Reader) (worldMap world.Map, err error) {
	worldMap = world.Map{}
	fileData, _ := ioutil.ReadAll(file)
	lines := strings.Split(string(fileData), "\n")
	for _, line := range lines {
		parseCity(line, worldMap)
	}
	return worldMap, nil
}
