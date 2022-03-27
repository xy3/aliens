package simulation

import (
	log "github.com/sirupsen/logrus"
	"github.com/xy3/aliens/world"
)

type Result struct {
	Days        int
	Exhausted   int
	Alive       int
	Dead        int
	Stuck       int
	WorldMap    world.Map
}

func (r Result) Display(logger *log.Logger) {
	logger.Info("==== SIMULATION RESULTS: ====")
	logger.WithFields(log.Fields{
		"daysPassed": r.Days,
		"dead":       r.Dead,
		"stuck":      r.Stuck,
		"exhausted":  r.Exhausted,
	}).Info("Result")
}
