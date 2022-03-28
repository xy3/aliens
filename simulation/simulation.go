package simulation

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/xy3/aliens/world"
)

type Simulation struct {
	aliens   []*world.Alien
	worldMap world.Map
	maxMoves uint
	events   chan string
	moves    chan world.Move
}

func New(aliens []*world.Alien, worldMap world.Map, maxMoves uint) Simulation {
	return Simulation{
		aliens:   aliens,
		worldMap: worldMap,
		maxMoves: maxMoves,
		events:   make(chan string),
		moves:    make(chan world.Move),
	}
}

func (s *Simulation) nextDay() (dead, stuck, exhausted int) {
	for _, alien := range s.aliens {
		if alien.Dead {
			dead++
			continue
		}
		if alien.Stuck {
			stuck++
			continue
		}
		if alien.Moves > s.maxMoves {
			exhausted++
			continue
		}
		move, err := alien.Move()
		if err != nil {
			log.Warn(err)
			alien.Dead = true
			continue
		}
		s.moves <- move
	}
	return
}

func (s *Simulation) DeployAliens() {
	for _, alien := range s.aliens {
		city, err := randomCity(s.worldMap)
		for err == nil && city.Destroyed {
			city, _ = randomCity(s.worldMap)
		}
		if err != nil {
			log.Warn("All cities have been destroyed during alien deployment")
			return
		}
		s.moves <- alien.DeployTo(city)
	}
}

// Run executes the Simulation for a given number of aliens and a path to a map file
func (s *Simulation) Run() Result {
	total := len(s.aliens)
	var dead, stuck, exhausted, days int
	for shouldContinue(total, dead, stuck, exhausted) {
		dead, stuck, exhausted = s.nextDay()
		days++
		if uint(days)%(s.maxMoves/10) == 0 {
			s.events <- fmt.Sprintf("%d DAYS HAVE PASSED", s.maxMoves/10)
		}
	}
	return Result{
		Days:      days,
		Exhausted: exhausted,
		Alive:     total - dead,
		Dead:      dead,
		Stuck:     stuck,
		WorldMap:  s.worldMap,
	}
}

func shouldContinue(total, dead, stuck, exhausted int) bool {
	return total > dead+stuck+exhausted
}
