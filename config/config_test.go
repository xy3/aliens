package config

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"path"
	"testing"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		file io.Reader
		wantErr bool
	}{
		{
			name:    "load invalid config",
			file:    bytes.NewBufferString(""),
			wantErr: true,
		},
		{
			name:    "load valid config",
			file:    bytes.NewBufferString("{\"MaxAlienMoves\": 100000, \"MapFile\": \"map.txt\", \"AlienNamesFile\": \"alien-names.txt\", \"DebugMode\": false}"),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Load(tt.file); (err != nil) != tt.wantErr {
				t.Errorf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPath(t *testing.T) {
	tests := []struct {
		name string
		want func () string
	}{
		{
			name: "basic config path",
			want: func() string {
				getwd, _ := os.Getwd()
				return path.Join(getwd, "config.json")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want(), Path())
		})
	}
}

func TestWrite(t *testing.T) {
	tests := []struct {
		name     string
		wantFile string
		wantErr  bool
	}{
		{
			name:     "default config write",
			wantFile: "{\n    \"MaxAlienMoves\": 100000,\n    \"MapFile\": \"map.txt\",\n    \"AlienNamesFile\": \"alien-names.txt\",\n    \"DebugMode\": false\n}",
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &bytes.Buffer{}
			err := Write(file)
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotFile := file.String()
			require.Equal(t, tt.wantFile, gotFile)
		})
	}
}
