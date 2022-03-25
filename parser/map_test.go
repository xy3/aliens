package parser

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"github.com/xy3/aliens/world"
	"io"
	"testing"
)

func TestParseMap(t *testing.T) {
	tests := []struct {
		name    string
		file    io.Reader
		want    world.Map
		wantErr bool
	}{
		{
			name:    "load empty map file",
			file:    bytes.NewBufferString(""),
			want:    world.Map{},
			wantErr: false,
		},
		{
			name: "load valid map file",
			file: bytes.NewBufferString("Baz"),
			want: world.Map{
				"Baz": &world.City{Name: "Baz"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseMap(tt.file)
			if (err != nil) != tt.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}
