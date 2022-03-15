package aliens

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"strings"
)

const AlienNamesFile = "alien-names.txt"

type Alien struct {
	Name  string
	City  *City
	Moves uint
	Dead  bool
}

func (a *Alien) GetMoves() uint {
	return a.Moves
}

func (a *Alien) Move(worldMap WorldMap) {
	// get available routes
	//if city == nil {
	//	log.Infof("%+v", a)
	//}
	routes := []*City{a.City.North, a.City.East, a.City.South, a.City.West}
	var availableRoutes []*City
	for _, r := range routes {
		if r != nil {
			availableRoutes = append(availableRoutes, r)
		}
	}
	if len(availableRoutes) == 0 {
		fmt.Println(a.Name, "is stuck with nowhere to go at", a.City.Name)
		a.Moves++
		return
	}
	// randomly choose a route
	chosenRoute := rand.Intn(len(availableRoutes))
	log.Info(len(availableRoutes), chosenRoute)
	newCity := availableRoutes[chosenRoute]
	log.WithField("alien", a.Name).WithField("city", a.City.Name).Infof("%s has moved to %s", a.Name, a.City.Name)

	if newCity.Inhabitant != nil && newCity.Inhabitant.Name != a.Name {
		log.WithFields(log.Fields{
			"opponents": a.Name + " vs " + newCity.Inhabitant.Name,
			"destroyedCity": a.City.Name,
		}).Infof("%s has been destroyed by %s and %s!", a.City.Name, a.Name, newCity.Inhabitant.Name)
		worldMap.DestroyCity(newCity)
		a.Dead = true
		a.City.Inhabitant.Dead = true
		return
	}
	// move along that route to the city
	a.City = newCity
	a.Moves++
	newCity.Inhabitant = a
	return
}

var alienNamesUsage map[string]int

func loadAlienNames(filePath string) ([]string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(file), "\n")
	alienNamesUsage = make(map[string]int, len(lines))
	return lines, nil
}

func randomName(names []string) string {
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

func randomCity(worldMap WorldMap) *City {
	uninhabitedWorldMap := WorldMap{}
	for _, city := range worldMap {
		if city.Inhabitant == nil {
			uninhabitedWorldMap[city.Name] = city
		}
	}
	usableMap := uninhabitedWorldMap
	if len(uninhabitedWorldMap) == 0 {
		usableMap = worldMap
	}
	randLocationInt := rand.Intn(len(usableMap))
	i := 0
	for _, city := range usableMap {
		if i == randLocationInt {
			return city
		}
		i++
	}
	return nil
}

func randomAliens(count int, worldMap WorldMap) ([]*Alien, error) {
	names, err := loadAlienNames(AlienNamesFile)
	if err != nil {
		return nil, err
	}
	alienList := make([]*Alien, count)
	for i := 0; i < count; i++ {
		city := randomCity(worldMap)
		alien := &Alien{Name: randomName(names), City: city}
		alien.City.Inhabitant = alien
		alienList[i] = alien
		log.Infof("%+v", alien)
	}
	return alienList, nil
}
