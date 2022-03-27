package simulation

import (
	"context"
	log "github.com/sirupsen/logrus"
	"github.com/xy3/aliens/world"
	"testing"
)

func TestSimulation_LogWorker(t *testing.T) {
	type fields struct {
		aliens   []*world.Alien
		worldMap world.Map
		maxMoves uint
		events   chan string
		moves    chan world.Move
	}
	type args struct {
		ctx    context.Context
		logger *log.Logger
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = &Simulation{
				aliens:   tt.fields.aliens,
				worldMap: tt.fields.worldMap,
				maxMoves: tt.fields.maxMoves,
				events:   tt.fields.events,
				moves:    tt.fields.moves,
			}
		})
	}
}

func Test_logMove(t *testing.T) {
	type args struct {
		move   world.Move
		logger *log.Logger
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}