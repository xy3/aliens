package aliens

import (
	"fmt"
	"io/ioutil"
	"strings"
)

var worldMap = WorldMap{}

type WorldMap map[string]*City

func (wm WorldMap) Serialize() (result string) {
	for _, c := range wm {
		if c != nil {
			result += c.String() + "\n"
		}
	}
	return strings.Trim(result, "\n")
}

// PrettyPrint formats the WorldMap into the original format using Serialize and prints the result
func (wm WorldMap) PrettyPrint() {
	fmt.Println(wm.Serialize())
}

// DestroyCity destroys a city and all of its connections to and from it
func (wm WorldMap) DestroyCity(city *City) {
	if city.North != nil {
		city.North.South = nil
	}
	if city.East != nil {
		city.East.West = nil
	}
	if city.West != nil {
		city.West.East = nil
	}
	if city.South != nil {
		city.South.North = nil
	}
	wm[city.Name] = nil
}


// loadMap reads in the map text file and parses the cities into the worldMap
func loadMap(filePath string) error {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	lines := strings.Split(string(file), "\n")
	for _, line := range lines {
		parseCity(line)
	}
	return nil
}

// parseCity reads in a line of input from the map file and updates cities in the worldMap
func parseCity(line string) {
	tokens := strings.Split(line, " ")

	newCityName := tokens[0]
	worldMap[newCityName] = &City{
		Name: newCityName,
	}
	newCity := worldMap[newCityName]

	for _, link := range tokens[1:] {
		parts := strings.Split(link, "=")
		direction, name := parts[0], parts[1]

		if worldMap[name] == nil {
			worldMap[name] = &City{Name: name}
		}
		linkCity := worldMap[name]

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
