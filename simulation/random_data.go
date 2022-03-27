package simulation

import (
	"errors"
	"fmt"
	"github.com/xy3/aliens/world"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

// randomCity tries to randomly select a city that does not have an inhabitant already from the worldMap
func randomCity(worldMap world.Map) (city *world.City, err error) {
	nonDestroyedMap := world.Map{}
	for _, city = range worldMap {
		if !city.Destroyed {
			nonDestroyedMap[city.Name] = city
		}
	}
	uninhabitedMap := world.Map{}
	for _, city = range worldMap {
		if !city.Destroyed && city.Inhabitant == nil {
			uninhabitedMap[city.Name] = city
		}
	}
	usableMap := uninhabitedMap
	if len(uninhabitedMap) == 0 {
		usableMap = nonDestroyedMap
	}
	if len(usableMap) == 0 {
		return nil, errors.New("all cities have been destroyed")
	}
	rand.Seed(time.Now().UnixNano())
	randLocationInt := rand.Intn(len(usableMap))
	// iterate through the usableMap until we are at the randomLocationInt position
	i := 0
	for _, city = range usableMap {
		if i == randLocationInt {
			return city, err
		}
		i++
	}
	return
}

// loadAlienNames reads in a text file of alien names one per line to be used when choosing random alien names
func loadAlienNames(file io.Reader) (names []string, err error) {
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	if len(fileData) > 1 {
		names = strings.Split(string(fileData), "\n")
		return
	}
	return []string{}, nil
}

// randomName generates a random number and uses that to select a random alien name from a provided list of names.
// It will append a number to the end of a name if that particular name has already been used for another alien.
func randomName(names []string, usage map[string]int) (string, error) {
	if len(names) == 0 {
		return "", errors.New("no names provided")
	}
	rand.Seed(time.Now().UnixNano())
	randNameInt := rand.Intn(len(names))
	name := names[randNameInt]
	if usage[name] == 0 {
		usage[name] = 1
		return name, nil
	}
	numberedName := fmt.Sprintf("%s_%d", name, usage[name])
	return numberedName, nil
}

// RandomAliens generates a list of Aliens with randomly selected names and initial city locations on the worldMap
func RandomAliens(count int, namesReader io.Reader) ([]*world.Alien, error) {
	names, err := loadAlienNames(namesReader)
	if err != nil {
		return nil, err
	}
	alienNamesUsage := map[string]int{}
	alienList := make([]*world.Alien, count)
	for i := 0; i < count; i++ {
		name, err := randomName(names, alienNamesUsage)
		if err != nil {
			return nil, err
		}
		alien := &world.Alien{
			Name: name,
		}
		alienList[i] = alien
		alienNamesUsage[name]++
	}
	return alienList, nil
}
