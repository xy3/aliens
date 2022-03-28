package world

import (
	log "github.com/sirupsen/logrus"
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
			tt.enemy.City.Inhabitant = tt.enemy
			tt.alien.fight(tt.enemy)
			require.True(t, tt.enemy.Dead)
			require.True(t, tt.alien.Dead)
			require.True(t, tt.enemy.City.Destroyed)
			require.False(t, tt.alien.City.Destroyed)
			require.Nil(t, tt.alien.City.Inhabitant)
		})
	}
}

func TestAlienMove(t *testing.T) {
	southCity := City{Name: "SouthCity"}.WithInhabitant(&Alien{Name: "EnemyAlien"}).WithNorth(&City{Name: "NorthCity"})

	tests := []struct {
		name      string
		alien     func() *Alien
		wantAlien *Alien
		wantErr   bool
		wantMove  Move
	}{
		{
			name: "alien is in non existing location",
			alien: func() *Alien {
				return &Alien{}
			},
			wantAlien: &Alien{},
			wantErr:   true,
		},
		{
			name: "alien is in destroyed location",
			alien: func() *Alien {
				return &Alien{City: &City{Destroyed: true}}
			},
			wantAlien: &Alien{City: &City{Destroyed: true}},
			wantErr:   true,
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
			wantMove: Move{
				Type: Stuck,
				City: City{Name: "Baz"},
				Alien: Alien{
					Name:  "alien",
					Dead:  false,
					City:  &City{Name: "Baz"},
					Stuck: true,
				},
			},
		},
		{
			name: "alien fights with inhabitant",
			alien: func() *Alien {
				return &Alien{Name: "alien", City: southCity.North}
			},
			wantAlien: &Alien{
				Name:  "alien",
				Dead:  true,
				City:  southCity.North,
				Stuck: false,
			},
			wantErr: false,
			wantMove: Move{
				Type:  Fight,
				City:  *southCity,
				Alien: Alien{Name: "alien", City: southCity.North},
				Enemy: *southCity.Inhabitant,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			alien := tt.alien()
			move, err := alien.Move()
			require.Equal(t, tt.wantErr, err != nil, "wantErr was %t but err was %t", tt.wantErr, err != nil)
			require.Equal(t, tt.wantMove, move)
			require.Equal(t, tt.wantAlien, alien)
		})
	}
}

func TestAlienMovesToOnlyExistingLocation(t *testing.T) {
	southCity := City{Name: "SouthCity"}.WithNorth(&City{Name: "NorthCity"})
	alien := &Alien{Name: "alien"}
	southCity.WithInhabitant(alien)

	wantAlien := &Alien{
		Name:  "alien",
		City:  southCity.North,
		Moves: 1,
	}

	t.Run("alien moves to only available next location", func(t *testing.T) {
		move, err := alien.Move()

		log.Infof("%+v", move)
		require.NoError(t, err)
		wantMove := Move{
			Type:  Moved,
			City:  *southCity.North.WithInhabitant(alien),
			Alien: *wantAlien,
		}
		require.Equal(t, wantMove, move)
		require.Equal(t, wantAlien, alien)
	})
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
				got := tt.alien.choosesToStay()
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

func TestAlienDeployTo(t *testing.T) {
	tests := []struct {
		name  string
		alien func() *Alien
		city  *City
		want  func() Move
	}{
		{
			name: "deploy to empty city",
			alien: func() *Alien {
				return &Alien{Name: "Roger"}
			},
			city: &City{Name: "Baz"},
			want: func() Move {
				baz := City{Name: "Baz"}.WithInhabitant(&Alien{Name: "Roger"})
				return Move{
					Type:  Moved,
					City:  *baz,
					Alien: *baz.Inhabitant,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.alien()
			got := a.DeployTo(tt.city)
			require.Equal(t, tt.want(), got)
		})
	}
}
