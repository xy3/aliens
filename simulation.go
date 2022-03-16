package aliens

import (
	log "github.com/sirupsen/logrus"
)

// RunSimulation executes the simulation for a given number of aliens and a path to a map file
func RunSimulation(alienCount int, mapFile string) error {
	err := loadMap(mapFile)
	if err != nil {
		return err
	}

	log.Info("PARSED SIMULATION WORLD MAP")
	worldMap.PrettyPrint()

	allAliens, err := randomAliens(alienCount, worldMap)
	if err != nil {
		return err
	}
	log.Infof("CREATED %d RANDOM ALIENS SUCCESSFULLY", alienCount)

	var aliensExceedMaxMoves bool
	var aliveAliens = alienCount
	for !aliensExceedMaxMoves && aliveAliens > 1 {
		aliveAliens = 0
		for i := 0; i < alienCount; i++ {
			alien := allAliens[i]
			if alien.Dead {
				continue
			}
			aliveAliens++
			alien.Move()
			aliensExceedMaxMoves = alien.Moves > Config.MaxAlienMoves
		}
	}

	log.Info("SIMULATION RESULTS:")
	if aliensExceedMaxMoves {
		log.Infof("All aliens have exceeded the max moves of %d", Config.MaxAlienMoves)
	}
	if aliveAliens == 0 {
		log.Info("All aliens have died")
	} else {
		log.Infof("%d alien remained and exhausted the move limit", aliveAliens)
	}

	log.Infof("The world map that still remains is:")
	worldMap.PrettyPrint()
	return nil
}
