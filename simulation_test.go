package aliens

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSimulationResultDisplay(t *testing.T) {
	tests := []struct {
		name   string
		result SimulationResult
	}{
		{
			name:   "displays results",
			result: SimulationResult{
				TotalAliens:          100,
				DaysPassed:           999,
				AliensExceedMaxMoves: false,
				AliveAliens:          0,
				StuckAliens:          0,
				WorldMapResult:       WorldMap{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.result.Display()
		})
	}
}

func TestSimulationResultShouldContinue(t *testing.T) {
	tests := []struct {
		name   string
		result SimulationResult
		want   bool
	}{
		{
			name: "should not continue if all aliens are dead",
			result: SimulationResult{
				AliensExceedMaxMoves: false,
				AliveAliens:          0,
				StuckAliens:          0,
			},
			want: false,
		},
		{
			name: "should not continue if all aliens exceed max moves",
			result: SimulationResult{
				AliensExceedMaxMoves: true,
				AliveAliens:          2,
				StuckAliens:          0,
			},
			want: false,
		},
		{
			name: "should not continue if all aliens are stuck",
			result: SimulationResult{
				AliensExceedMaxMoves: false,
				AliveAliens:          2,
				StuckAliens:          2,
			},
			want: false,
		},
		{
			name: "should not continue if one alien is left alive",
			result: SimulationResult{
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

func TestRandomCity(t *testing.T) {
	tests := []struct {
		name     string
		worldMap WorldMap
		wantCity *City
	}{
		{
			name:     "randomly selects a city",
			worldMap: WorldMap{"Baz": &City{Name: "Baz"}},
			wantCity: &City{Name: "Baz"},
		},
		{
			name:     "returns nil for empty map",
			worldMap: WorldMap{},
			wantCity: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCity := randomCity(tt.worldMap)
			require.Equal(t, tt.wantCity, gotCity)
		})
	}
}

func TestRandomName(t *testing.T) {
	alienNamesUsage["usedName"] = 1

	tests := []struct {
		name  string
		names []string
		want  string
	}{
		{
			name:  "randomly selects a name from a list",
			names: []string{"testname"},
			want:  "testname",
		},
		{
			name:  "appends number to a used name",
			names: []string{"usedName"},
			want:  "usedName_1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := randomName(tt.names)
			require.Equal(t, tt.want, got)
		})
	}
}
