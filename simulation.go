package aliens

import log "github.com/sirupsen/logrus"

const MaxAlienMoves = 100000

func RunSimulation(alienCount int, mapFile string) error {
	worldMap, err := loadMap(mapFile)
	if err != nil {
		return err
	}

	log.Info("SIMULATION WORLD MAP")
	worldMap.PrettyPrint()

	allAliens, err := randomAliens(alienCount, worldMap)
	if err != nil {
		return err
	}
	log.Infof("CREATED %d RANDOM ALIENS SUCCESSFULLY", alienCount)

	var aliensExceedMaxMoves bool
	var aliensAreDead bool
	for !aliensExceedMaxMoves && !aliensAreDead {
		for i := 0; i < alienCount; i++ {
			alien := allAliens[i]
			if alien.Dead {
				aliensAreDead = true
				continue
			}
			aliensAreDead = false
			alien.Move(worldMap)
			aliensExceedMaxMoves = alien.Moves > MaxAlienMoves
		}
	}

	log.Infof("aliensExceedMaxMoves: %t, aliensAreDead: %t", aliensExceedMaxMoves, aliensAreDead)
	return nil
}

