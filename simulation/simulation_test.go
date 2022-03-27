package simulation

import (
	"github.com/stretchr/testify/require"
	"github.com/xy3/aliens/world"
	"reflect"
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
	type fields struct {
		aliens   []*world.Alien
		worldMap world.Map
		maxMoves uint
		events   chan string
		moves    chan world.Move
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//s := &Simulation{
			//	aliens:   tt.fields.aliens,
			//	worldMap: tt.fields.worldMap,
			//	maxMoves: tt.fields.maxMoves,
			//	events:   tt.fields.events,
			//	moves:    tt.fields.moves,
			//}
		})
	}
}

func TestSimulation_Run(t *testing.T) {
	type fields struct {
		aliens   []*world.Alien
		worldMap world.Map
		maxMoves uint
		events   chan string
		moves    chan world.Move
	}
	tests := []struct {
		name   string
		fields fields
		want   Result
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Simulation{
				aliens:   tt.fields.aliens,
				worldMap: tt.fields.worldMap,
				maxMoves: tt.fields.maxMoves,
				events:   tt.fields.events,
				moves:    tt.fields.moves,
			}
			if got := s.Run(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimulation_nextDay(t *testing.T) {
	type fields struct {
		aliens   []*world.Alien
		worldMap world.Map
		maxMoves uint
		events   chan string
		moves    chan world.Move
	}
	tests := []struct {
		name          string
		fields        fields
		wantDead      int
		wantStuck     int
		wantExhausted int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Simulation{
				aliens:   tt.fields.aliens,
				worldMap: tt.fields.worldMap,
				maxMoves: tt.fields.maxMoves,
				events:   tt.fields.events,
				moves:    tt.fields.moves,
			}
			gotDead, gotStuck, gotExhausted := s.nextDay()
			if gotDead != tt.wantDead {
				t.Errorf("nextDay() gotDead = %v, want %v", gotDead, tt.wantDead)
			}
			if gotStuck != tt.wantStuck {
				t.Errorf("nextDay() gotStuck = %v, want %v", gotStuck, tt.wantStuck)
			}
			if gotExhausted != tt.wantExhausted {
				t.Errorf("nextDay() gotExhausted = %v, want %v", gotExhausted, tt.wantExhausted)
			}
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := shouldContinue(tt.args.total, tt.args.dead, tt.args.stuck, tt.args.exhausted); got != tt.want {
				t.Errorf("shouldContinue() = %v, want %v", got, tt.want)
			}
		})
	}
}