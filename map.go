package aliens

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

type WorldMap map[string]*City

func (wm WorldMap) Serialize() (result string) {
	var cities []string
    for _, city := range wm {
		if city != nil {
			cities = append(cities, city.String())
		}
    }
	sort.Strings(cities)
	return strings.Join(cities, "\n")
}

// PrettyPrint formats the WorldMap into the original format using Serialize and prints the result
func (wm WorldMap) PrettyPrint() {
	fmt.Println(wm.Serialize())
}

// LoadMap reads in the map text file and parses the cities into the worldMap
func LoadMap(filePath string) (worldMap WorldMap, err error) {
	file, err := fs.Open(filePath)
	if err != nil {
		return nil, err
	}
	worldMap = WorldMap{}
	fileData, _ := ioutil.ReadAll(file)
	lines := strings.Split(string(fileData), "\n")
	for _, line := range lines {
		parseCity(line, worldMap)
	}
	return worldMap, nil
}

// parseCity reads in a line of input from the map file and updates cities in the worldMap
func parseCity(line string, worldMap WorldMap) {
	if len(line) == 0 {
		return
	}
	tokens := strings.Split(line, " ")

	newCityName := tokens[0]
	worldMap[newCityName] = &City{
		Name: newCityName,
	}
	newCity := worldMap[newCityName]

	for _, link := range tokens[1:] {
		parts := strings.Split(link, "=")
		direction, linkedCityName := parts[0], parts[1]

		if worldMap[linkedCityName] == nil {
			worldMap[linkedCityName] = &City{Name: linkedCityName}
		}
		linkCity := worldMap[linkedCityName]

		switch strings.ToLower(direction) {
		case "north":
			newCity.North = linkCity
			linkCity.South = newCity
			break
		case "east":
			newCity.East = linkCity
			linkCity.West = newCity
			break
		case "south":
			newCity.South = linkCity
			linkCity.North = newCity
			break
		case "west":
			newCity.West = linkCity
			linkCity.East = newCity
			break
		}
	}
}
