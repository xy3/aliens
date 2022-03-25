package simulation

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xy3/aliens/world"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

// randomCity tries to randomly select a city that does not have an inhabitant already from the worldMap
func randomCity(worldMap world.Map) (city *world.City) {
	uninhabitedWorldMap := world.Map{}
	for _, city = range worldMap {
		if city.Inhabitant == nil {
			uninhabitedWorldMap[city.Name] = city
		}
	}
	usableMap := uninhabitedWorldMap
	if len(uninhabitedWorldMap) == 0 {
		usableMap = worldMap
	}
	if len(usableMap) == 0 {
		return
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

var alienNamesUsage = map[string]int{}

// loadAlienNames reads in a text file of Alien names one per line to be used when choosing random Alien names
func loadAlienNames(file io.Reader) ([]string, error) {
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(fileData), "\n")
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

// RandomAliens generates a list of Aliens with randomly selected names and initial city locations on the worldMap
func RandomAliens(count int, worldMap world.Map, namesReader io.Reader) ([]*world.Alien, error) {
	names, err := loadAlienNames(namesReader)
	if err != nil {
		return nil, err
	}
	alienList := make([]*world.Alien, count)
	for i := 0; i < count; i++ {
		city := randomCity(worldMap)
		alien := &world.Alien{Name: randomName(names), City: city}
		alien.City.Inhabitant = alien
		alienList[i] = alien
		log.Debugf("Created alien: %+v", *alien)
	}
	return alienList, nil
}