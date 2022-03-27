package simulation

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"github.com/xy3/aliens/world"
	"io"
	"testing"
)

func TestRandomAliens(t *testing.T) {
	type args struct {
		count       int
		namesReader io.Reader
	}
	tests := []struct {
		name    string
		args    args
		want    []*world.Alien
		wantErr bool
	}{
		{
			name: "no aliens",
			args: args{
				count:       0,
				namesReader: bytes.NewBufferString(""),
			},
			want:    []*world.Alien{},
			wantErr: false,
		},
		{
			name: "10 aliens with no names provided",
			args: args{
				count:       10,
				namesReader: bytes.NewBufferString(""),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "1 alien with 1 name provided",
			args: args{
				count:       1,
				namesReader: bytes.NewBufferString("Name1"),
			},
			want:    []*world.Alien{{Name: "Name1"}},
			wantErr: false,
		},
		{
			name: "2 aliens with 1 name provided",
			args: args{
				count:       2,
				namesReader: bytes.NewBufferString("Name1"),
			},
			want:    []*world.Alien{{Name: "Name1"}, {Name: "Name1_2"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RandomAliens(tt.args.count, tt.args.namesReader)
			require.Equal(t, tt.wantErr, err != nil, "wantErr was %t but err was %t", tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}

func Test_loadAlienNames(t *testing.T) {
	tests := []struct {
		name    string
		file    io.Reader
		want    []string
		wantErr bool
	}{
		{
			name:    "load two names",
			file:    bytes.NewBufferString("FirstName\nSecondName"),
			want:    []string{"FirstName", "SecondName"},
			wantErr: false,
		},
		{
			name:    "load one name",
			file:    bytes.NewBufferString("FirstName"),
			want:    []string{"FirstName"},
			wantErr: false,
		},
		{
			name:    "load no names",
			file:    bytes.NewBufferString(""),
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loadAlienNames(tt.file)
			if tt.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

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
			gotCity, _ := randomCity(tt.worldMap)
			require.Equal(t, tt.wantCity, gotCity)
		})
	}
}

func TestRandomName(t *testing.T) {
	tests := []struct {
		name    string
		names   []string
		want    string
		wantErr bool
		usage   map[string]int
	}{
		{
			name:  "randomly selects a name from a list",
			names: []string{"testname"},
			want:  "testname",
			usage: map[string]int{},
		},
		{
			name:  "appends number to a used name",
			names: []string{"usedName"},
			want:  "usedName_1",
			usage: map[string]int{"usedName": 1},
		},
		{
			name:    "empty name list",
			names:   []string{},
			want:    "",
			usage:   map[string]int{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := randomName(tt.names, tt.usage)
			require.Equal(t, tt.wantErr, err != nil, "wantErr was %t but err was %t", tt.wantErr, err != nil)
			require.Equal(t, tt.want, got)
		})
	}
}
