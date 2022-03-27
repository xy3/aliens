package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCitySerialize(t *testing.T) {
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
			got := tt.city.serialize()
			require.Equal(t, tt.want, got)
		})
	}
}
