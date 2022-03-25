package simulation

import (
	log "github.com/sirupsen/logrus"
	"github.com/xy3/aliens/world"
)

type Simulation struct {
	aliens   []*world.Alien
	worldMap world.Map
	maxMoves uint
	result   Result
	done     chan bool
	DoneSimulation     chan bool
	events   chan string
	moves    chan world.Move
}

func (s *Simulation) nextDay() {
	s.result.DaysPassed++
	s.result.AliveAliens = 0
	s.result.StuckAliens = 0
	for i := 0; i < len(s.aliens); i++ {
		currentAlien := s.aliens[i]
		if currentAlien.Dead {
			continue
		}
		s.result.AliveAliens++
		if currentAlien.Stuck {
			s.result.StuckAliens++
			continue
		}

		currentAlien.Move(s.worldMap, s.moves)
		s.result.AliensExceedMaxMoves = currentAlien.Moves > s.maxMoves
	}
	if s.result.DaysPassed%500 == 0 {
		s.events <- "500 DAYS HAVE PASSED"
	}
}

func New(aliens []*world.Alien, worldMap world.Map, maxMoves uint) Simulation {
	return Simulation{
		aliens:   aliens,
		worldMap: worldMap,
		maxMoves: maxMoves,
		result: Result{
			TotalAliens: len(aliens),
			AliveAliens: len(aliens),
			WorldMap:    worldMap,
			DaysPassed:  1,
		},
		done:  make(chan bool, 1),
		DoneSimulation:  make(chan bool, 1),
		moves: make(chan world.Move, 100),
	}
}

// Run executes the Simulation for a given number of aliens and a path to a map file
func (s *Simulation) Run() {
	for s.result.shouldContinue() {
		s.nextDay()
	}
	s.result.WorldMap = s.worldMap
	s.done <- true
}

func (s Simulation) Result() Result {
	return s.result
}

func (s *Simulation) LogWorker(logger *log.Logger) {
	for {
		select {
		case msg := <-s.events:
			logger.Info(msg)
		case move := <-s.moves:
			logMove(move, logger)
		case <-s.done:
			s.DoneSimulation <- true
			return
		}
	}
}

func logMove(move world.Move, logger *log.Logger) {
	fields := log.Fields{
		"alien": move.AlienName,
		"city":  move.City,
	}

	switch move.MoveType {
	case world.Moved:
		logger.WithFields(fields).Debugf("%s has moved to %s", move.AlienName, move.City)
	case world.Stays:
		logger.WithFields(fields).Debugf("%s decides to remain at %s", move.AlienName, move.City)
	case world.Stuck:
		logger.WithFields(fields).Debugf("%s is stuck at %s", move.AlienName, move.City)
	case world.Fight:
		fields = log.Fields{
			"opponents":     move.AlienName + " vs " + move.EnemyName,
			"destroyedCity": move.City,
		}
		logger.WithFields(fields).Infof("%s has been destroyed by %s and %s", move.City, move.AlienName, move.EnemyName)
	}
}
