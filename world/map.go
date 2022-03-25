package world

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

type Map map[string]*City

func (m Map) Serialize() (result string) {
	var cities []string
    for _, city := range m {
		if city != nil {
			cities = append(cities, city.serialize())
		}
    }
	sort.Strings(cities)
	return strings.Join(cities, "\n")
}

// Print formats the WorldMap into the original format using Serialize and prints the result
func (m Map) Print(w io.Writer) {
	fmt.Fprintln(w, m.Serialize())
}