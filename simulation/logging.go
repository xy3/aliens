package simulation

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/xy3/aliens/world"
)

func (s *Simulation) LogWorker(ctx context.Context, logger *log.Logger) {
	for {
		select {
		case msg := <-s.events:
			logger.Info(msg)
		case move := <-s.moves:
			logMove(move, logger)
		case <-ctx.Done():
			return
		}
	}
}

func logMove(move world.Move, logger *log.Logger) {
	fields := log.Fields{
		"alien": move.Alien.Name,
		"city":  move.City.Name,
	}

	switch move.Type {
	case world.Moved:
		logger.WithFields(fields).Debugf("%s has moved to %s", move.Alien.Name, move.City.Name)
	case world.Stays:
		logger.WithFields(fields).Debugf("%s decides to remain at %s", move.Alien.Name, move.City.Name)
	case world.Stuck:
		logger.WithFields(fields).Debugf("%s is stuck at %s", move.Alien.Name, move.City.Name)
	case world.Fight:
		fields = log.Fields{
			"opponents":     move.Alien.Name + " vs " + move.Enemy.Name,
			"destroyedCity": move.City.Name,
		}
		logger.WithFields(fields).Infof("%s has been destroyed by %s and %s", move.City.Name, move.Alien.Name, move.Enemy.Name)
	}
}
