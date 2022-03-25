package world

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMapSerialize(t *testing.T) {
	tests := []struct {
		name       string
		wm         Map
		wantResult string
	}{
		{
			name:       "serialize empty map",
			wm:         nil,
			wantResult: "",
		},
		{
			name:       "serialize map with one city",
			wm:         Map{"Baz": &City{Name: "Baz"}},
			wantResult: "Baz",
		},
		{
			name: "serialize map with two cities",
			wm: Map{
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
