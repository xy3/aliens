package world

import (
	"math/rand"
	"time"
)

// Alien represents an Alien in the simulation. It has a Name, a Moves counter, a bool to mark it's alive state and a
// City to record its present location.
type Alien struct {
	Name  string
	City  *City
	Moves uint
	Dead  bool
	Stuck bool
}

type MoveType int

const (
	Stays MoveType = iota
	Stuck
	Moved
	Fight
)

type Move struct {
	MoveType  MoveType
	City      string
	AlienName string
	EnemyName string
}

// Move moves an alien to a next location by randomly choosing between staying in the same location or moving to another
// city that has a connection to its current city. If there are no paths available the alien will be forced to stay at
// the same location and update its move counter.
func (a *Alien) Move(worldMap Map, moves chan Move) {
	// prevent aliens from accidentally existing in non-existing cities
	if a.City == nil || worldMap[a.City.Name] == nil {
		a.Dead = true
		return
	}

	// if there is an inhabitant (and it isn't the current alien) in the current city, start a fight between them
	if a.City.Inhabitant != nil && a.City.Inhabitant.Name != a.Name {
		moves <- Move{
			MoveType:  Fight,
			City:      a.City.Name,
			AlienName: a.Name,
			EnemyName: a.City.Inhabitant.Name,
		}
		a.fight(a.City.Inhabitant, worldMap)
		return
	}

	// alien randomly chooses not to move to a new location
	if a.choosesToStay() {
		moves <- Move{
			MoveType:  Stays,
			City:      a.City.Name,
			AlienName: a.Name,
		}
		return
	}

	// find routes the alien can move to, or stay at the same location
	newCity := a.findAvailableRoute()
	if newCity == nil {
		moves <- Move{
			MoveType:  Stuck,
			City:      a.City.Name,
			AlienName: a.Name,
		}
		return
	}

	// if there is an inhabitant (and it isn't the current alien), start a fight between them
	if newCity.Inhabitant != nil && newCity.Inhabitant.Name != a.Name {
		newCity.Inhabitant.City = newCity
		moves <- Move{
			MoveType:  Fight,
			City:      newCity.Name,
			AlienName: a.Name,
			EnemyName: newCity.Inhabitant.Name,
		}
		a.fight(newCity.Inhabitant, worldMap)
		return
	}

	// move this alien along that route to the next city
	a.City = newCity
	a.Moves++
	newCity.Inhabitant = a
	moves <- Move{
		MoveType:  Moved,
		City:      a.City.Name,
		AlienName: a.Name,
	}
}

// findAvailableRoute randomly chooses an available route the Alien can move to
func (a *Alien) findAvailableRoute() *City {
	// filter out the nil routes (routes that lead nowhere)
	routes := []*City{a.City.North, a.City.East, a.City.South, a.City.West}
	var availableRoutes []*City
	for _, r := range routes {
		if r != nil {
			availableRoutes = append(availableRoutes, r)
		}
	}
	if len(availableRoutes) == 0 {
		// here the alien is stuck at the current location and cannot move anywhere else
		a.Stuck = true
		return nil
	}
	// randomly choose a route to a next city
	chosenRoute := rand.Intn(len(availableRoutes))
	return availableRoutes[chosenRoute]
}

// choosesToStay gives Alien a 1 in 15 chance to stay at the same location
func (a *Alien) choosesToStay() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(15) == 1
}

// fight sets both of the Alien's Dead property to true and destroys the city they are in, including its connections
func (a *Alien) fight(enemy *Alien, worldMap Map) {
	a.City = nil
	enemy.City.Inhabitant.Dead = true
	enemy.City.Destroy(worldMap)
	enemy.City = nil
	a.Dead = true
	return
}
