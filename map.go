package aliens

import (
	"fmt"
	"io/ioutil"
	"strings"
)

type WorldMap map[string]*City

func (wm WorldMap) Serialize() (result string) {
	for _, c := range wm {
		if c != nil {
			result += c.String() + "\n"
		}
	}
	return strings.Trim(result, "\n")
}

func (wm WorldMap) PrettyPrint() {
	fmt.Println(wm.Serialize())
}

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

type City struct {
	Name       string
	North      *City
	East       *City
	South      *City
	West       *City
	Inhabitant *Alien
}

func (c *City) String() string {
	var routes string
	if c.North != nil {
		routes += "north=" + c.North.Name + " "
	}
	if c.East != nil {
		routes += "east=" + c.East.Name + " "
	}
	if c.South != nil {
		routes += "south=" + c.South.Name + " "
	}
	if c.West != nil {
		routes += "west=" + c.West.Name + " "
	}
	return fmt.Sprintf("%s %v", c.Name, strings.Trim(routes, " "))
}

// loadMap reads in the map text file and parses the cities
func loadMap(filePath string) (WorldMap, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(file), "\n")
	worldMap := make(WorldMap, len(lines))
	for _, line := range lines {
		parseCity(line, worldMap)
	}
	return worldMap, nil
}

// parseCity reads in a line of input from the map file and updates
// cities in the worldMap
func parseCity(line string, worldMap WorldMap) {
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

		switch direction {
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
