package parser

import (
	"github.com/xy3/aliens/world"
	"strings"
)

// parseCity reads in a line of input from the map file and updates cities in the worldMap
func parseCity(line string, worldMap world.Map) {
	if len(line) == 0 {
		return
	}
	tokens := strings.Split(line, " ")

	newCityName := tokens[0]
	worldMap[newCityName] = &world.City{
		Name: newCityName,
	}
	newCity := worldMap[newCityName]

	for _, link := range tokens[1:] {
		parts := strings.Split(link, "=")
		direction, linkedCityName := parts[0], parts[1]

		if worldMap[linkedCityName] == nil {
			worldMap[linkedCityName] = &world.City{Name: linkedCityName}
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
