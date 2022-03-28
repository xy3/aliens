package world

import (
	"fmt"
	"strings"
)

// City represents a city in the worldMap. It has a Name and a pointer to each of the Cities in the cardinal directions.
// If there is no city in a certain direction the pointer will be null. Inhabitant records the presently occupying alien
// of the City.
type City struct {
	Name       string
	North      *City
	East       *City
	South      *City
	West       *City
	Inhabitant *Alien
	Destroyed  bool
}

func (c City) WithNorth(north *City) *City {
	c.North = north
	north.South = &c
	return &c
}

func (c City) WithEast(east *City) *City {
	c.East = east
	east.West = &c
	return &c
}

func (c City) WithSouth(south *City) *City {
	c.South = south
	south.North = &c
	return &c
}

func (c City) WithWest(west *City) *City {
	c.West = west
	west.East = &c
	return &c
}

func (c City) WithInhabitant(inhabitant *Alien) *City {
	c.Inhabitant = inhabitant
	inhabitant.City = &c
	return &c
}

// serialize returns a representation of the City used when formatting the worldMap at the end of the simulation.
func (c *City) serialize() string {
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
	return strings.Trim(fmt.Sprintf("%s %v", c.Name, strings.Trim(routes, " ")), " ")
}
