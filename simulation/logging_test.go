package simulation

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"github.com/xy3/aliens/world"
	"sync"
	"testing"
)

func TestSimulationLogWorker(t *testing.T) {
	t.Run("log worker logs events and moves", func(t *testing.T) {
		s := Simulation{
			events: make(chan string),
			moves:  make(chan world.Move),
		}

		logger, hook := test.NewNullLogger()
		logger.SetLevel(log.DebugLevel)
		ctx, cancel := context.WithCancel(context.Background())
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.LogWorker(ctx, logger)
		}()

		s.moves <- world.Move{
			Type:  world.Moved,
			City:  world.City{Name: "testCity"},
			Alien: world.Alien{Name: "Test"},
		}
		s.moves <- world.Move{
			Type:  world.Fight,
			City:  world.City{Name: "testCity"},
			Alien: world.Alien{Name: "Test"},
			Enemy: world.Alien{Name: "TestEnemy"},
		}
		s.events <- "test event"
		wantMessages := 3
		cancel()
		wg.Wait()
		require.Equal(t, wantMessages, len(hook.Entries))
		hook.Reset()
		require.Nil(t, hook.LastEntry())
	})
}

func TestLogMove(t *testing.T) {
	t.Run("log fight move message", func(t *testing.T) {
		move := world.Move{
			Type:  world.Fight,
			City:  world.City{Name: "testCity"},
			Alien: world.Alien{Name: "Test"},
			Enemy: world.Alien{Name: "TestEnemy"},
		}
		logger, hook := test.NewNullLogger()
		logger.SetLevel(log.DebugLevel)
		logMove(move, logger)
		wantMessages := 1
		require.Equal(t, wantMessages, len(hook.Entries))
		hook.Reset()
		require.Nil(t, hook.LastEntry())
	})
}
