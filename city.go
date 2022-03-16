package aliens

import (
	"fmt"
	"strings"
)

// City represents a city in the worldMap. It has a Name and a pointer to each of the Cities in the cardinal directions.
// If there is no city in a certain direction the pointer will be null. Inhabitant records the presently occupying Alien
// of the City.
type City struct {
	Name       string
	North      *City
	East       *City
	South      *City
	West       *City
	Inhabitant *Alien
}

// String returns a representation of the City used when formatting the worldMap at the end of the simulation.
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
