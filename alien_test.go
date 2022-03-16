package aliens

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAlienFight(t *testing.T) {
	tests := []struct {
		name  string
		alien *Alien
		enemy *Alien
	}{
		{
			name: "fight over a city",
			alien: &Alien{
				Name: "alien1",
				City: &City{Name: "Baz"},
				Dead: false,
			},
			enemy: &Alien{
				Name: "alien2",
				City: &City{Name: "Bar"},
				Dead: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			worldMap := WorldMap{}
			tt.enemy.City.Inhabitant = tt.enemy
			var city = *tt.enemy.City
			tt.alien.Fight(tt.enemy, worldMap)
			require.True(t, tt.enemy.Dead)
			require.True(t, tt.alien.Dead)
			require.Nil(t, tt.enemy.City)
			require.Nil(t, worldMap[city.Name])
		})
	}
}

func TestAlienMove(t *testing.T) {
	// some test cities
	northCity := &City{Name: "NorthCity"}
	southCity := &City{Name: "SouthCity", North: northCity}
	northCity.South = southCity

	tests := []struct {
		name      string
		alien     func() *Alien
		worldMap  func() WorldMap
		wantAlien *Alien
	}{
		{
			name: "alien moves to non existing location",
			alien: func() *Alien {
				return &Alien{}
			},
			wantAlien: &Alien{Dead: true},
			worldMap: func() WorldMap {
				return WorldMap{}
			},
		},
		{
			name: "alien is forced to remain at current location",
			alien: func() *Alien {
				return &Alien{Name: "alien", City: &City{Name: "Baz"}}
			},
			wantAlien: &Alien{
				Name:  "alien",
				Dead:  false,
				City:  &City{Name: "Baz"},
				Stuck: true,
			},
			worldMap: func() WorldMap {
				baz := &City{Name: "Baz"}
				return WorldMap{baz.Name: baz}
			},
		},
		{
			name: "alien fights at current location",
			alien: func() *Alien {
				baz := &City{Name: "Baz"}
				baz.Inhabitant = &Alien{Name: "enemy", City: baz}
				return &Alien{Name: "alien", City: baz}
			},
			wantAlien: &Alien{Name: "alien", Dead: true, City: nil},
			worldMap: func() WorldMap {
				baz := &City{Name: "Baz"}
				return WorldMap{baz.Name: baz}
			},
		},
		{
			name: "alien moves to only available next location",
			alien: func() *Alien {
				return &Alien{Name: "alien", City: southCity}
			},
			wantAlien: &Alien{
				Name:  "alien",
				City:  northCity,
				Moves: 1,
			},
			worldMap: func() WorldMap {
				return WorldMap{
					southCity.Name: southCity,
					northCity.Name: northCity,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alien := tt.alien()
			alien.Move(tt.worldMap())
			require.Equal(t, tt.wantAlien, alien)
		})
	}
}

func TestAlienChoosesToRemain(t *testing.T) {
	tests := []struct {
		name  string
		alien *Alien
	}{
		{
			name: "alien chooses to remain randomly",
			alien: &Alien{
				Name: "Bob",
				City: &City{Name: "Baz"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var trueTimes int
			for i := 0; i < 100; i++ {
				got := tt.alien.choosesToRemain()
				if got {
					trueTimes++
				}
			}
			require.NotEqual(t, 0, trueTimes)
			require.NotEqual(t, 100, trueTimes)
		})
	}
}

func TestAlienFindAvailableRoute(t *testing.T) {
	tests := []struct {
		name      string
		alien     *Alien
		want      *City
		wantStuck bool
	}{
		{
			name: "alien is stuck with no availableRoutes",
			alien: &Alien{
				Name: "alien1",
				City: &City{Name: "Baz"},
			},
			want:      nil,
			wantStuck: true,
		},
		{
			name: "alien chooses only available route",
			alien: &Alien{
				Name: "alien2",
				City: &City{Name: "Baz", North: &City{Name: "Foo"}},
			},
			want:      &City{Name: "Foo"},
			wantStuck: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.alien.findAvailableRoute()
			require.Equal(t, tt.want, got)
			require.Equal(t, tt.wantStuck, tt.alien.Stuck)
		})
	}
}

func TestAlienMoveTo(t *testing.T) {
	tests := []struct {
		name  string
		alien *Alien
		city  *City
	}{
		{
			name:  "alien moves to city",
			alien: &Alien{},
			city:  &City{Name: "Baz"},
		},
		{
			name:  "alien moves to same city",
			alien: &Alien{City: &City{Name: "Baz"}},
			city:  &City{Name: "Baz"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wantMoves := tt.alien.Moves + 1
			tt.alien.moveTo(tt.city)
			require.Equal(t, tt.city, tt.alien.City)
			require.Equal(t, wantMoves, tt.alien.Moves)
			require.Equal(t, tt.alien, tt.city.Inhabitant)
		})
	}
}
