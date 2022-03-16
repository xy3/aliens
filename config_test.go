package aliens

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name       string
		writeFiles func()
		wantErr    bool
	}{
		{
			name: "fails to read non existing config file",
			writeFiles: func() {
				fs = afero.NewMemMapFs()
			},
			wantErr: true,
		},
		{
			name: "fails to read non valid config file",
			writeFiles: func() {
				fs = afero.NewMemMapFs()
				_ = afero.WriteFile(fs, ConfigFile, []byte("non-valid"), 0644)
			},
			wantErr: true,
		},
				{
			name: "reads valid config file",
			writeFiles: func() {
				fs = afero.NewMemMapFs()
				_ = afero.WriteFile(
					fs,
					ConfigFile,
					[]byte(`{"MaxAlienMoves":10,"MapFile":"map.txt","AlienNamesFile":"alien-names.txt"}`),
					0644,
				)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.writeFiles()
			err := LoadConfig()
			if tt.wantErr {
				require.Error(t, err)
			}
		})
	}
}

func TestWriteConfig(t *testing.T) {
	fs = afero.NewMemMapFs()
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "writes config correctly",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WriteConfig()
			if tt.wantErr {
				require.Error(t, err)
			}
		})
	}
}
