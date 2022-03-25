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
				TotalAliens:          100,
				DaysPassed:           999,
				AliensExceedMaxMoves: false,
				AliveAliens:          0,
				StuckAliens:          0,
				WorldMap:             world.Map{},
			},
			wantMessages: []log.Entry{
				{Message: "==== SIMULATION RESULTS: ====", Data: log.Fields{}},
				{Message: "SimulationResult", Data: log.Fields{"days passed": 999, "dead": 100, "exhausted": 0, "trapped": 0}},
				{Message: "999 days have passed in the Simulation", Data: log.Fields{}},
				{Message: "100 aliens have died", Data: log.Fields{}},
				{Message: "0 aliens exhausted the move limit", Data: log.Fields{}},
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

func TestResultShouldContinue(t *testing.T) {
	tests := []struct {
		name   string
		result Result
		want   bool
	}{
		{
			name: "should not continue if all aliens are dead",
			result: Result{
				AliensExceedMaxMoves: false,
				AliveAliens:          0,
				StuckAliens:          0,
			},
			want: false,
		},
		{
			name: "should not continue if all aliens exceed max moves",
			result: Result{
				AliensExceedMaxMoves: true,
				AliveAliens:          2,
				StuckAliens:          0,
			},
			want: false,
		},
		{
			name: "should not continue if all aliens are stuck",
			result: Result{
				AliensExceedMaxMoves: false,
				AliveAliens:          2,
				StuckAliens:          2,
			},
			want: false,
		},
		{
			name: "should not continue if one alien is left alive",
			result: Result{
				AliensExceedMaxMoves: false,
				AliveAliens:          1,
				StuckAliens:          0,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.result.shouldContinue()
			require.Equal(t, tt.want, got)
		})
	}
}
