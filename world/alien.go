package world

import (
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

// Move moves an Alien to a next location by randomly choosing between staying in the same location or moving to another
// city that has a connection to its current city. If there are no paths available the Alien will be forced to stay at
// the same location and update its move counter.
func (a *Alien) Move(worldMap Map) {
	// prevent aliens from accidentally existing in non-existing cities
	if a.City == nil || worldMap[a.City.Name] == nil {
		a.Dead = true
		return
	}

	// if there is an inhabitant (and it isn't the current Alien) in the current city, start a fight between them
	if a.City.Inhabitant != nil && a.City.Inhabitant.Name != a.Name {
		a.Fight(a.City.Inhabitant, worldMap)
		return
	}

	// alien randomly chooses not to move to a new location
	if a.choosesToRemain() {
		return
	}

	// find routes the alien can move to, or stay at the same location
	newCity := a.findAvailableRoute()
	if newCity == nil {
		return
	}

	// if there is an inhabitant (and it isn't the current Alien), start a fight between them
	if newCity.Inhabitant != nil && newCity.Inhabitant.Name != a.Name {
		newCity.Inhabitant.City = newCity
		a.Fight(newCity.Inhabitant, worldMap)
		return
	}

	// move this alien along that route to the next city
	a.moveTo(newCity)
	return
}

// moveTo moves an Alien to a new city by updating the Alien's city, it's Moves and the City's Inhabitant
func (a *Alien) moveTo(city *City) {
	a.City = city
	a.Moves++
	city.Inhabitant = a

	log.WithFields(log.Fields{
		"alien": a.Name,
		"city":  a.City.Name,
	}).Debugf("%s has moved to %s", a.Name, a.City.Name)
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
		log.WithFields(log.Fields{
			"alien": a.Name,
			"city":  a.City.Name,
		}).Debugf("%s is stuck at %s and cannot move", a.Name, a.City.Name)
		a.Stuck = true
		return nil
	}
	// randomly choose a route to a next city
	chosenRoute := rand.Intn(len(availableRoutes))
	return availableRoutes[chosenRoute]
}

// choosesToRemain gives Alien a 1 in 15 chance to stay at the same location
func (a *Alien) choosesToRemain() bool {
	rand.Seed(time.Now().UnixNano())
	remains := rand.Intn(15) == 1
	if remains {
		log.WithFields(log.Fields{
			"alien": a.Name,
			"city":  a.City.Name,
		}).Debugf("%s decides to remain at %s", a.Name, a.City.Name)
	}
	return remains
}

// Fight sets both of the Alien's Dead property to true and destroys the city they are in, including its connections
func (a *Alien) Fight(enemy *Alien, worldMap Map) {
	log.WithFields(log.Fields{
		"opponents":     a.Name + " vs " + enemy.Name,
		"destroyedCity": enemy.City.Name,
	}).Infof("%s has been destroyed by %s and %s!", enemy.City.Name, a.Name, enemy.Name)
	a.City = nil
	enemy.City.Inhabitant.Dead = true
	enemy.City.Destroy(worldMap)
	enemy.City = nil
	a.Dead = true
	return
}
