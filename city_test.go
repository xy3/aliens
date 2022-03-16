package aliens

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCityDestroy(t *testing.T) {
	tests := []struct {
		name string
		city func() *City
	}{
		{
			name: "city with all connections destroys correctly",
			city: func() *City {
				baz := &City{
					Name:       "Baz",
					North:      &City{Name: "NorthCity"},
					East:       &City{Name: "EastCity"},
					South:      &City{Name: "SouthCity"},
					West:       &City{Name: "WestCity"},
				}
				baz.North.South = baz
				baz.East.West = baz
				baz.South.North = baz
				baz.West.East = baz
				return baz
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			city := tt.city()
			cityName := city.Name
			worldMap := WorldMap{
				city.Name: city,
			}
			city.Destroy(worldMap)
			require.Nil(t, city.North.South)
			require.Nil(t, city.East.West)
			require.Nil(t, city.South.North)
			require.Nil(t, city.West.East)
			require.Nil(t, worldMap[cityName])
		})
	}
}

func TestCityString(t *testing.T) {
	tests := []struct {
		name string
		city *City
		want string
	}{
		{
			name: "city with no connections",
			city: &City{Name: "Baz"},
			want: "Baz",
		},
		{
			name: "city with all connections",
			city: &City{
				Name:       "Baz",
				North:      &City{Name: "NorthCity"},
				East:       &City{Name: "EastCity"},
				South:      &City{Name: "SouthCity"},
				West:       &City{Name: "WestCity"},
			},
			want: "Baz north=NorthCity east=EastCity south=SouthCity west=WestCity",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.city.String()
			require.Equal(t, tt.want, got)
		})
	}
}
