package aliens

import (
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

func TestLoadMap(t *testing.T) {
	tests := []struct {
		name       string
		filePath   string
		want       WorldMap
		wantErr    bool
		writeFiles func()
	}{
		{
			name:     "load empty map file",
			filePath: "map.txt",
			want:     nil,
			wantErr:  true,
			writeFiles: func() {
				fs = afero.NewMemMapFs()
			},
		},
		{
			name:     "load valid map file",
			filePath: "map.txt",
			want: WorldMap{
				"Baz": &City{Name: "Baz"},
			},
			wantErr: false,
			writeFiles: func() {
				fs = afero.NewMemMapFs()
				_ = afero.WriteFile(fs, "map.txt", []byte("Baz"), 0644)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.writeFiles()
			got, err := LoadMap(tt.filePath)
			if (err != nil) != tt.wantErr {
				require.Error(t, err)
			}
			require.Equal(t, tt.want, got)
		})
	}
}

func TestWorldMapPrettyPrint(t *testing.T) {
	tests := []struct {
		name        string
		wm          WorldMap
		expectedLog string
	}{
		{
			name:        "pretty print empty map",
			wm:          nil,
			expectedLog: "\n",
		},
		{
			name: "pretty print map with one city",
			wm: WorldMap{
				"Baz": &City{Name: "Baz"},
			},
			expectedLog: "Baz\n",
		},
		{
			name: "pretty print map with two cities",
			wm: WorldMap{
				"Baz": &City{Name: "Baz"},
				"Foo": &City{Name: "Foo", South: &City{Name: "Bar"}},
			},
			expectedLog: "Baz\nFoo south=Bar\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			tt.wm.PrettyPrint()
			w.Close()
			out, _ := ioutil.ReadAll(r)
			os.Stdout = oldStdout
			require.Equal(t, tt.expectedLog, string(out))
		})
	}
}

func TestWorldMapSerialize(t *testing.T) {
	tests := []struct {
		name       string
		wm         WorldMap
		wantResult string
	}{
		{
			name:       "serialize empty map",
			wm:         nil,
			wantResult: "",
		},
		{
			name: "serialize map with one city",
			wm: WorldMap{
				"Baz": &City{Name: "Baz"},
			},
			wantResult: "Baz",
		},
		{
			name: "serialize map with two cities",
			wm: WorldMap{
				"Baz": &City{Name: "Baz"},
				"Foo": &City{Name: "Foo", South: &City{Name: "Bar"}},
			},
			wantResult: "Baz\nFoo south=Bar",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult := tt.wm.Serialize()
			require.Equal(t, tt.wantResult, gotResult)
		})
	}
}

func TestParseCity(t *testing.T) {
	tests := []struct {
		name             string
		line             string
		expectedWorldMap func() WorldMap
	}{
		{
			name: "parse blank line",
			line: "",
			expectedWorldMap: func() WorldMap {
				return WorldMap{}
			},
		},
		{
			name: "parse city with no connections",
			line: "Baz",
			expectedWorldMap: func() WorldMap {
				baz := &City{Name: "Baz"}
				return WorldMap{baz.Name: baz}
			},
		},
		{
			name: "parse city with one connection",
			line: "Baz north=Foo",
			expectedWorldMap: func() WorldMap {
				baz := &City{Name: "Baz"}
				foo := &City{Name: "Foo"}
				baz.North = foo
				foo.South = baz
				return WorldMap{
					baz.Name: baz,
					foo.Name: foo,
				}
			},
		},
		{
			name: "parse city with all connections",
			line: "Baz north=Foo east=Bar south=Foobar west=Vegas",
			expectedWorldMap: func() WorldMap {
				baz := &City{
					Name:  "Baz",
					North: &City{Name: "Foo"},
					East:  &City{Name: "Bar"},
					South: &City{Name: "Foobar"},
					West:  &City{Name: "Vegas"},
				}
				baz.North.South = baz
				baz.East.West = baz
				baz.South.North = baz
				baz.West.East = baz
				return WorldMap{
					baz.Name:       baz,
					baz.North.Name: baz.North,
					baz.East.Name:  baz.East,
					baz.South.Name: baz.South,
					baz.West.Name:  baz.West,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testWorldMap := WorldMap{}
			parseCity(tt.line, testWorldMap)
			require.Equal(t, tt.expectedWorldMap(), testWorldMap)
		})
	}
}
