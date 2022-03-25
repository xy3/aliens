package simulation

import (
	"github.com/stretchr/testify/require"
	"github.com/xy3/aliens/world"
	"testing"
)

func TestRandomCity(t *testing.T) {
	tests := []struct {
		name     string
		worldMap world.Map
		wantCity *world.City
	}{
		{
			name:     "randomly selects a city",
			worldMap: world.Map{"Baz": &world.City{Name: "Baz"}},
			wantCity: &world.City{Name: "Baz"},
		},
		{
			name:     "returns nil for empty map",
			worldMap: world.Map{},
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
