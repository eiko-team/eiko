package reseachhelper

import (
	"github.com/eiko-team/eiko/misc/structures"
)

// ConsumableToSearchable Extract data from structures.Consumable to build a
// new struct with searchable values
func ConsumableToSearchable(c structures.Consumable) interface{} {
	return struct {
		Name       string
		ID         int64
		Company    string
		Categories []string
		Ingredient []string
		Label      []string
	}{
		Name:       c.Name,
		ID:         c.ID,
		Company:    c.Company,
		Categories: c.Categories,
		Ingredient: c.Ingredient,
		Label:      c.Label,
	}
}
