package aliens

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

// Alien represents an Alien in the simulation. It has a Name, a Moves counter, a bool to mark it's alive state and a
// City to record its present location.
type Alien struct {
	Name  string
	City  *City
	Moves uint
	Dead  bool
}

// Move moves an Alien to a next location by randomly choosing between staying in the same location or moving to another
// city that has a connection to its current city. If there are no paths available the Alien will be forced to stay at
// the same location and update its move counter.
func (a *Alien) Move() {
	// give aliens a 1 in 15 chance of staying at the same location
	rand.Seed(time.Now().UnixNano())
	remainsStill := rand.Intn(15) == 1
	if remainsStill {
		log.WithFields(log.Fields{
			"alien": a.Name,
			"city":  a.City.Name,
		}).Debugf("%s decides to remain at %s", a.Name, a.City.Name)
		return
	}

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
		a.Moves++
		return
	}
	// randomly choose a route to a next city
	chosenRoute := rand.Intn(len(availableRoutes))
	newCity := availableRoutes[chosenRoute]
	log.WithFields(log.Fields{
		"alien": a.Name,
		"city":  a.City.Name,
	}).Debugf("%s has moved to %s", a.Name, a.City.Name)

	// if there is an inhabitant (and it isn't the current Alien), start a fight between them
	if newCity.Inhabitant != nil && newCity.Inhabitant.Name != a.Name {
		a.Fight(newCity.Inhabitant)
		return
	}
	// move this alien along that route to the next city
	a.City = newCity
	a.Moves++
	newCity.Inhabitant = a
	return
}

// Fight sets both of the Alien's Dead property to true and destroys the city they are in, including its connections
func (a *Alien) Fight(enemy *Alien) {
	log.WithFields(log.Fields{
		"opponents":      a.Name + " vs " + enemy.Name,
		"destroyed city": a.City.Name,
	}).Infof("%s has been destroyed by %s and %s!", a.City.Name, a.Name, enemy.Name)
	worldMap.DestroyCity(enemy.City)
	a.Dead = true
	a.City.Inhabitant.Dead = true
	return
}

var alienNamesUsage map[string]int

// loadAlienNames reads in a text file of Alien names one per line to be used when choosing random Alien names
func loadAlienNames(filePath string) ([]string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(file), "\n")
	alienNamesUsage = make(map[string]int, len(lines))
	return lines, nil
}

// randomName generates a random number and uses that to select a random alien name from a provided list of names.
// It will append a number to the end of a name if that particular name has already been used for another Alien.
func randomName(names []string) string {
	rand.Seed(time.Now().UnixNano())
	randNameInt := rand.Intn(len(names))
	name := names[randNameInt]
	if alienNamesUsage[name] == 0 {
		alienNamesUsage[name] = 1
		return name
	}
	numberedName := fmt.Sprintf("%s_%d", name, alienNamesUsage[name])
	alienNamesUsage[name]++
	return numberedName
}

// randomCity tries to randomly select a city that does not have an inhabitant already from the worldMap
func randomCity(worldMap WorldMap) (city *City) {
	uninhabitedWorldMap := WorldMap{}
	for _, city = range worldMap {
		if city.Inhabitant == nil {
			uninhabitedWorldMap[city.Name] = city
		}
	}
	usableMap := uninhabitedWorldMap
	if len(uninhabitedWorldMap) == 0 {
		usableMap = worldMap
	}
	rand.Seed(time.Now().UnixNano())
	randLocationInt := rand.Intn(len(usableMap))
	// iterate through the usableMap until we are at the randomLocationInt position
	i := 0
	for _, city = range usableMap {
		if i == randLocationInt {
			return city
		}
		i++
	}
	return
}

// randomAliens generates a list of Aliens with randomly selected names and initial city locations on the worldMap
func randomAliens(count int, worldMap WorldMap) ([]*Alien, error) {
	names, err := loadAlienNames(Config.AlienNamesFile)
	if err != nil {
		return nil, err
	}
	alienList := make([]*Alien, count)
	for i := 0; i < count; i++ {
		city := randomCity(worldMap)
		alien := &Alien{Name: randomName(names), City: city}
		alien.City.Inhabitant = alien
		alienList[i] = alien
		log.Debugf("Created alien: %+v", *alien)
	}
	return alienList, nil
}
