package aliens

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

// RunSimulation executes the simulation for a given number of aliens and a path to a map file
func RunSimulation(alienCount int, worldMap WorldMap) (*SimulationResult, error) {
	allAliens, err := randomAliens(alienCount, worldMap)
	if err != nil {
		return nil, err
	}
	log.Infof("CREATED %d RANDOM ALIENS SUCCESSFULLY", alienCount)
	result := &SimulationResult{
		TotalAliens:    alienCount,
		AliveAliens:    alienCount,
		WorldMapResult: worldMap,
		DaysPassed:     1,
	}

	for result.shouldContinue() {
		result.DaysPassed++
		result.AliveAliens = 0
		result.StuckAliens = 0
		for i := 0; i < alienCount; i++ {
			alien := allAliens[i]
			if alien.Dead {
				continue
			}
			result.AliveAliens++
			if alien.Stuck {
				result.StuckAliens++
				continue
			}
			alien.Move(worldMap)
			result.AliensExceedMaxMoves = alien.Moves > Config.MaxAlienMoves
		}
		if result.DaysPassed%500 == 0 {
			log.Info("500 DAYS HAVE PASSED")
		}
	}
	result.WorldMapResult = worldMap
	return result, nil
}

type SimulationResult struct {
	TotalAliens          int
	DaysPassed           int
	AliensExceedMaxMoves bool
	AliveAliens          int
	StuckAliens          int
	WorldMapResult       WorldMap
}

func (r SimulationResult) shouldContinue() bool {
	return !r.AliensExceedMaxMoves && r.AliveAliens > 1 && r.StuckAliens < r.AliveAliens
}

func (r SimulationResult) Display() {
	log.Info("==== SIMULATION RESULTS: ====")
	log.WithFields(log.Fields{
		"days passed": r.DaysPassed,
		"dead": r.TotalAliens-r.AliveAliens,
		"trapped": r.StuckAliens,
		"exhausted": r.AliveAliens-r.StuckAliens,
	}).Info("SimulationResult")
	log.Infof("%d days have passed in the simulation", r.DaysPassed)
	if r.AliensExceedMaxMoves {
		log.Infof("All aliens have exceeded the max moves of %d", Config.MaxAlienMoves)
	}
	log.Infof("%d aliens have died", r.TotalAliens-r.AliveAliens)
	if r.StuckAliens > 0 {
		log.Infof("%d aliens got trapped in cities", r.StuckAliens)
	}
	log.Infof("%d aliens exhausted the move limit", r.AliveAliens-r.StuckAliens)

	log.Infof("The world map that still remains is:")
	r.WorldMapResult.PrettyPrint()
}

var alienNamesUsage = map[string]int{}

// loadAlienNames reads in a text file of Alien names one per line to be used when choosing random Alien names
func loadAlienNames(filePath string) ([]string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(file), "\n")
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
