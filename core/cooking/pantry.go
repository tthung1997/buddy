package cooking

import (
	"errors"
	"time"
)

var ErrPantryItemNotFound = errors.New("pantry item not found")

type PantryItem struct {
	IngredientId    string    `json:"ingredientId"`
	Amount          float64   `json:"amount"`
	Unit            Unit      `json:"unit"`
	UpdatedDateTime time.Time `json:"updatedDateTime"`
}

type IPantryRepository interface {
	List() ([]PantryItem, error)
	Get(ingredientId string) (PantryItem, error)
	CreateOrUpdate(PantryItem) error
	// ReplaceAll swaps the whole pantry in a single step so a bulk edit cannot
	// interleave with another writer.
	ReplaceAll([]PantryItem) error
	Delete(ingredientId string) error
}
