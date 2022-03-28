package simulation

import (
	"github.com/stretchr/testify/require"
	"github.com/xy3/aliens/world"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		aliens   []*world.Alien
		worldMap world.Map
		maxMoves uint
	}
	tests := []struct {
		name string
		args args
		want Simulation
	}{
		{
			name: "nil arguments",
			args: args{
				aliens:   nil,
				worldMap: nil,
				maxMoves: 0,
			},
			want: Simulation{
				aliens:   nil,
				worldMap: nil,
				maxMoves: 0,
			},
		},
		{
			name: "typical arguments",
			args: args{
				aliens:   []*world.Alien{},
				worldMap: world.Map{},
				maxMoves: 100,
			},
			want: Simulation{
				aliens:   []*world.Alien{},
				worldMap: world.Map{},
				maxMoves: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.aliens, tt.args.worldMap, tt.args.maxMoves)
			require.Equal(t, tt.want.aliens, got.aliens)
			require.Equal(t, tt.want.worldMap, got.worldMap)
			require.Equal(t, tt.want.maxMoves, got.maxMoves)
			require.NotNil(t, got.moves)
			require.NotNil(t, got.events)
		})
	}
}

func TestSimulation_DeployAliens(t *testing.T) {
	tests := []struct {
		name       string
		simulation Simulation
		wantMove   func() world.Move
	}{
		{
			name: "deploy one alien",
			simulation: Simulation{
				aliens:   []*world.Alien{{Name: "Baz"}},
				worldMap: world.Map{"Bar": &world.City{Name: "Bar"}},
				maxMoves: 10,
				events:   nil,
				moves:    make(chan world.Move, 1),
			},
			wantMove: func() world.Move {
				city := world.City{Name: "Bar"}.WithInhabitant(&world.Alien{Name: "Baz"})
				return world.Move{
					Type:  world.Moved,
					City:  *city,
					Alien: *city.Inhabitant,
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.simulation.DeployAliens()
			require.Equal(t, tt.wantMove(), <-tt.simulation.moves)
		})
	}
}

func TestSimulationRun(t *testing.T) {
	tests := []struct {
		name       string
		simulation Simulation
		want       Result
	}{
		{
			name:       "simulation with no aliens",
			simulation: Simulation{},
			want:       Result{},
		},
		{
			name: "simulation with two aliens",
			simulation: Simulation{
				aliens:   []*world.Alien{{Name: "Roger", City: &world.City{Name: "Baz"}}, {Name: "Foo", City: &world.City{Name: "Bar"}}},
				worldMap: nil,
				maxMoves: 200,
				moves:    make(chan world.Move, 10),
			},
			want: Result{Exhausted: 0, Alive: 2, Dead: 0, Stuck: 2, WorldMap: world.Map(nil)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.simulation.Run()
			got.Days = 0
			require.Equal(t, tt.want, got)
		})
	}
}

func TestSimulation_nextDay(t *testing.T) {
	tests := []struct {
		name          string
		simulation    Simulation
		wantDead      int
		wantStuck     int
		wantExhausted int
	}{
		{
			name: "simulation with no aliens",
			simulation: Simulation{
				aliens: nil,
			},
			wantDead:      0,
			wantStuck:     0,
			wantExhausted: 0,
		},
		{
			name: "alien that exceeds max moves",
			simulation: Simulation{
				aliens:   []*world.Alien{{Name: "Roger", Moves: 10000}},
				worldMap: world.Map{"Baz": &world.City{Name: "Baz"}},
				maxMoves: 1,
				events:   nil,
				moves:    nil,
			},
			wantDead:      0,
			wantStuck:     0,
			wantExhausted: 1,
		},
		{
			name: "alien that is stuck",
			simulation: Simulation{
				aliens:   []*world.Alien{{Name: "Roger", Moves: 1, Stuck: true}},
				maxMoves: 100,
			},
			wantDead:      0,
			wantStuck:     1,
			wantExhausted: 0,
		},
		{
			name: "alien that is dead",
			simulation: Simulation{
				aliens:   []*world.Alien{{Name: "Roger", Moves: 1, Dead: true}},
				maxMoves: 100,
			},
			wantDead:      1,
			wantStuck:     0,
			wantExhausted: 0,
		},
		{
			name: "alien that can move",
			simulation: Simulation{
				aliens:   []*world.Alien{{Name: "Roger", Moves: 1}},
				maxMoves: 100,
				moves:    make(chan world.Move, 1),
			},
			wantDead:      0,
			wantStuck:     0,
			wantExhausted: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDead, gotStuck, gotExhausted := tt.simulation.nextDay()
			require.Equal(t, tt.wantDead, gotDead)
			require.Equal(t, tt.wantStuck, gotStuck)
			require.Equal(t, tt.wantExhausted, gotExhausted)
		})
	}
}

func Test_shouldContinue(t *testing.T) {
	type args struct {
		total     int
		dead      int
		stuck     int
		exhausted int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "basic should continue when total aliens is not greater than dead+stuck+exhausted",
			args: args{
				total:     10,
				dead:      100,
				stuck:     23,
				exhausted: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldContinue(tt.args.total, tt.args.dead, tt.args.stuck, tt.args.exhausted); got != tt.want {
				t.Errorf("shouldContinue() = %v, want %v", got, tt.want)
			}
		})
	}
}
