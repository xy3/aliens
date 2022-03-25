package world

import (
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