package simulation

import (
	log "github.com/sirupsen/logrus"
	"github.com/xy3/aliens/config"
	"github.com/xy3/aliens/world"
)

type Result struct {
	TotalAliens          int
	DaysPassed           int
	AliensExceedMaxMoves bool
	AliveAliens          int
	StuckAliens          int
	WorldMapResult       world.Map
}

func (r Result) shouldContinue() bool {
	return !r.AliensExceedMaxMoves && r.AliveAliens > 1 && r.StuckAliens < r.AliveAliens
}

func (r Result) Display(logger *log.Logger) {
	logger.Info("==== SIMULATION RESULTS: ====")
	logger.WithFields(log.Fields{
		"days passed": r.DaysPassed,
		"dead":        r.TotalAliens - r.AliveAliens,
		"trapped":     r.StuckAliens,
		"exhausted":   r.AliveAliens - r.StuckAliens,
	}).Info("SimulationResult")
	logger.Infof("%d days have passed in the Simulation", r.DaysPassed)
	if r.AliensExceedMaxMoves {
		logger.Infof("All aliens have exceeded the max moves of %d", config.Config.MaxAlienMoves)
	}
	logger.Infof("%d aliens have died", r.TotalAliens-r.AliveAliens)
	if r.StuckAliens > 0 {
		logger.Infof("%d aliens got trapped in cities", r.StuckAliens)
	}
	logger.Infof("%d aliens exhausted the move limit", r.AliveAliens-r.StuckAliens)
}
