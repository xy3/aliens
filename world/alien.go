package world

import (
	"errors"
	log "github.com/sirupsen/logrus"
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
	Type  MoveType
	City  City
	Alien Alien
	Enemy Alien
}

// Move moves an alien to a next location by randomly choosing between staying in the same location or moving to another
// city that has a connection to its current city. If there are no paths available the alien will be forced to stay at
// the same location and update its move counter.
func (a *Alien) Move() (Move, error) {
	if a.City == nil || a.City.Destroyed {
		return Move{}, errors.New("alien is currently in a non existing or destroyed city")
	}
	// alien randomly chooses not to move to a new location
	if a.choosesToStay() {
		return Move{
			Type:  Stays,
			City:  *a.City,
			Alien: *a,
		}, nil
	}

	// find routes the alien can move to, or stay at the same location
	newCity := a.findAvailableRoute()
	if newCity == nil {
		return Move{
			Type:  Stuck,
			City:  *a.City,
			Alien: *a,
		}, nil
	}

	// if there is an inhabitant (and it isn't the current alien), start a fight between them
	if newCity.Inhabitant != nil && newCity.Inhabitant != a {
		move := Move{
			Type:  Fight,
			City:  *newCity,
			Alien: *a,
			Enemy: *newCity.Inhabitant,
		}
		a.fight(newCity.Inhabitant)
		return move, nil
	}

	// move this alien along that route to the next city
	return a.moveTo(newCity), nil
}

func (a *Alien) moveTo(newCity *City) Move {
	a.City.Inhabitant = nil
	newCity.Inhabitant = a
	a.City = newCity
	a.Moves++
	return Move{
		Type:  Moved,
		City:  *a.City,
		Alien: *a,
	}
}

// findAvailableRoute randomly chooses an available route the Alien can move to
func (a *Alien) findAvailableRoute() *City {
	// filter out the nil routes (routes that lead nowhere)
	routes := []*City{a.City.North, a.City.East, a.City.South, a.City.West}
	var availableRoutes []*City
	for _, r := range routes {
		if r != nil && !r.Destroyed {
			availableRoutes = append(availableRoutes, r)
		}
	}
	if len(availableRoutes) == 0 {
		// here the alien is stuck at the current location and cannot move anywhere else
		a.Stuck = true
		return nil
	}
	// randomly choose a route to a next city
	rand.Seed(time.Now().UnixNano())
	chosenRoute := rand.Intn(len(availableRoutes))
	return availableRoutes[chosenRoute]
}

// choosesToStay gives Alien a 1 in 15 chance to stay at the same location
func (a *Alien) choosesToStay() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(15) == 1
}

// fight sets both of the Alien's Dead property to true and destroys the city they are in, including its connections
func (a *Alien) fight(defender *Alien) {
	a.City.Inhabitant = nil
	defender.City.Destroyed = true
	a.Dead = true
	defender.Dead = true
}

func (a *Alien) DeployTo(city *City) Move {
	a.City = city
	if a.City.Inhabitant != nil && a.City.Inhabitant != a {
		move := Move{
			Type:  Fight,
			City:  *a.City,
			Alien: *a,
			Enemy: *a.City.Inhabitant,
		}
		a.fight(a.City.Inhabitant)
		return move
	}
	log.Infof("Deployed %s to %s", a.Name, a.City.Name)
	a.City.Inhabitant = a
	return Move{
		Type:  Moved,
		City:  *a.City,
		Alien: *a,
	}
}
