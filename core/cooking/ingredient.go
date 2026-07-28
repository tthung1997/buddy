package cooking

import (
	"errors"
	"time"
)

var ErrIngredientNotFound = errors.New("ingredient not found")

type Ingredient struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Category     string   `json:"category,omitempty"`
	ImageURL     string   `json:"imageUrl,omitempty"`
	DefaultUnit  Unit     `json:"defaultUnit"`
	AllowedUnits []Unit   `json:"allowedUnits,omitempty"`
	Tags         []string `json:"tags,omitempty"`
	// Staple marks an ingredient nobody measures what is left of. Staples are
	// tracked by presence through PantryItem.Stock rather than by amount, and
	// cooking never deducts them.
	Staple bool `json:"staple,omitempty"`
	// GramsPerUnit is the mass of one DefaultUnit of this ingredient. It stays
	// zero when unknown and exists so density aware volume to mass conversion
	// can be added later without reshaping stored data.
	GramsPerUnit    float64   `json:"gramsPerUnit,omitempty"`
	UpdatedDateTime time.Time `json:"updatedDateTime"`
}

// Allows reports whether the ingredient may be measured in the given unit. An
// ingredient without an explicit allow list accepts any unit of the same kind
// as its default unit.
func (i Ingredient) Allows(unit Unit) bool {
	if !unit.IsValid() {
		return false
	}
	if len(i.AllowedUnits) == 0 {
		if !i.DefaultUnit.IsValid() {
			return true
		}
		defaultKind, err := i.DefaultUnit.Kind()
		if err != nil {
			return false
		}
		kind, err := unit.Kind()
		if err != nil {
			return false
		}
		return kind == defaultKind
	}
	for _, allowed := range i.AllowedUnits {
		if allowed == unit {
			return true
		}
	}
	return unit == i.DefaultUnit
}

type IIngredientRepository interface {
	List() ([]Ingredient, error)
	Get(id string) (Ingredient, error)
	CreateOrUpdate(Ingredient) error
	// ReplaceAll swaps the whole catalog in a single step so a bulk edit cannot
	// interleave with another writer.
	ReplaceAll([]Ingredient) error
	Delete(id string) error
}
