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
}

func (s *Simulation) NextDay() {
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
		currentAlien.Move(s.worldMap)
		s.result.AliensExceedMaxMoves = currentAlien.Moves > s.maxMoves
	}
	if s.result.DaysPassed%500 == 0 {
		log.Info("500 DAYS HAVE PASSED")
	}
}

func New(aliens []*world.Alien, worldMap world.Map, maxMoves uint) Simulation {
	return Simulation{
		aliens:   aliens,
		worldMap: worldMap,
		result: Result{
			TotalAliens:    len(aliens),
			AliveAliens:    len(aliens),
			WorldMapResult: worldMap,
			DaysPassed:     1,
		},
		maxMoves: maxMoves,
	}
}

// Run executes the Simulation for a given number of aliens and a path to a map file
func (s *Simulation) Run() {
	for s.result.shouldContinue() {
		s.NextDay()
	}
	s.result.WorldMapResult = s.worldMap
}
