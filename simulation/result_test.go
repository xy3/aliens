package simulation

import (
	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/require"
	"github.com/xy3/aliens/world"
	"testing"
)

func TestResultDisplay(t *testing.T) {
	tests := []struct {
		name         string
		result       Result
		wantMessages []log.Entry
	}{
		{
			name: "displays results",
			result: Result{
				Days:      999,
				Exhausted: 12,
				Alive:     0,
				Dead:      100,
				Stuck:     2,
				WorldMap:  world.Map{},
			},
			wantMessages: []log.Entry{
				{Message: "==== SIMULATION RESULTS: ====", Data: log.Fields{}},
				{Message: "Result", Data: log.Fields{"daysPassed": 999, "dead": 100, "exhausted": 12, "stuck": 2}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, hook := test.NewNullLogger()
			tt.result.Display(logger)
			require.Equal(t, len(tt.wantMessages), len(hook.Entries))
			for i, msg := range hook.Entries {
				require.Equal(t, tt.wantMessages[i].Message, msg.Message)
				require.Equal(t, tt.wantMessages[i].Data, msg.Data)
			}

			hook.Reset()
			require.Nil(t, hook.LastEntry())
		})
	}
}
