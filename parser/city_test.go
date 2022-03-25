package parser

import (
	"github.com/stretchr/testify/require"
	"github.com/xy3/aliens/world"
	"testing"
)


func TestParseCity(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		expectedWorldMap func() world.Map
	}{
		{
			name: "parse blank line",
			line: "",
			expectedWorldMap: func() world.Map {
				return world.Map{}
			},
		},
		{
			name: "parse city with no connections",
			line: "Baz",
			expectedWorldMap: func() world.Map {
				baz := &world.City{Name: "Baz"}
				return world.Map{baz.Name: baz}
			},
		},
		{
			name: "parse city with one connection",
			line: "Baz north=Foo",
			expectedWorldMap: func() world.Map {
				baz := &world.City{Name: "Baz"}
				foo := &world.City{Name: "Foo"}
				baz.North = foo
				foo.South = baz
				return world.Map{
					baz.Name: baz,
					foo.Name: foo,
				}
			},
		},
		{
			name: "parse city with all connections",
			line: "Baz north=Foo east=Bar south=Foobar west=Vegas",
			expectedWorldMap: func() world.Map {
				baz := &world.City{
					Name:  "Baz",
					North: &world.City{Name: "Foo"},
					East:  &world.City{Name: "Bar"},
					South: &world.City{Name: "Foobar"},
					West:  &world.City{Name: "Vegas"},
				}
				baz.North.South = baz
				baz.East.West = baz
				baz.South.North = baz
				baz.West.East = baz
				return world.Map{
					baz.Name:       baz,
					baz.North.Name: baz.North,
					baz.East.Name:  baz.East,
					baz.South.Name: baz.South,
					baz.West.Name:  baz.West,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWorldMap := world.Map{}
			parseCity(tt.line, testWorldMap)
			require.Equal(t, tt.expectedWorldMap(), testWorldMap)
		})
	}
}