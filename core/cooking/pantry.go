package cooking

import (
	"errors"
	"strings"
	"time"
)

var ErrPantryItemNotFound = errors.New("pantry item not found")

// StockLevel describes a staple that is tracked by presence rather than by
// amount. Salt, oil, and rice are the kind of thing nobody measures what is
// left of; they are either there, running out, or gone.
type StockLevel string

const (
	// StockUnknown is the zero value and means the entry is tracked by amount.
	StockUnknown StockLevel = ""
	// InStock means there is enough to cook with.
	InStock StockLevel = "IN_STOCK"
	// LowStock still cooks tonight but belongs on the next shopping list.
	LowStock StockLevel = "LOW"
	// OutOfStock blocks any recipe that needs the ingredient.
	OutOfStock StockLevel = "OUT"
)

var orderedStockLevels = []StockLevel{InStock, LowStock, OutOfStock}

// StockLevels returns the selectable stock levels, healthiest first.
func StockLevels() []StockLevel {
	levels := make([]StockLevel, len(orderedStockLevels))
	copy(levels, orderedStockLevels)
	return levels
}

func (s StockLevel) IsValid() bool {
	return s == InStock || s == LowStock || s == OutOfStock
}

// Available reports whether there is enough of the staple to cook with.
func (s StockLevel) Available() bool {
	return s == InStock || s == LowStock
}

func (s StockLevel) String() string {
	return string(s)
}

// Label renders the level for display.
func (s StockLevel) Label() string {
	switch s {
	case InStock:
		return "In stock"
	case LowStock:
		return "Running low"
	case OutOfStock:
		return "Out of stock"
	default:
		return ""
	}
}

// ParseStockLevel resolves user supplied text to a stock level, defaulting to
// in stock so a staple row is never accidentally treated as empty.
func ParseStockLevel(value string) StockLevel {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case string(LowStock), "LOW_STOCK", "RUNNING LOW", "RUNNING_LOW":
		return LowStock
	case string(OutOfStock), "OUT_OF_STOCK", "OUT OF STOCK", "NONE":
		return OutOfStock
	default:
		return InStock
	}
}

type PantryItem struct {
	IngredientId string  `json:"ingredientId"`
	Amount       float64 `json:"amount,omitempty"`
	Unit         Unit    `json:"unit,omitempty"`
	// Stock is used instead of Amount and Unit when the ingredient is a staple.
	Stock           StockLevel `json:"stock,omitempty"`
	UpdatedDateTime time.Time  `json:"updatedDateTime"`
}

// Pantry is what is in the kitchen together with the catalog that describes it.
// The catalog is what tells the engine which entries are staples, so the two
// always travel together.
type Pantry struct {
	Items   []PantryItem
	Catalog map[string]Ingredient
}

// NewPantry bundles pantry items with the catalog describing them.
func NewPantry(items []PantryItem, catalog map[string]Ingredient) Pantry {
	return Pantry{Items: items, Catalog: catalog}
}

// IsStaple reports whether an ingredient is tracked by presence.
func (p Pantry) IsStaple(ingredientId string) bool {
	return p.Catalog[ingredientId].Staple
}

// StockOf returns the stock level recorded for a staple. A staple with no
// pantry entry at all counts as out of stock.
func (p Pantry) StockOf(ingredientId string) StockLevel {
	for _, item := range p.Items {
		if item.IngredientId != ingredientId {
			continue
		}
		if item.Stock.IsValid() {
			return item.Stock
		}
		// A staple saved before it was marked as one carries an amount instead
		// of a level; treat any amount at all as being in stock.
		if item.Amount > 0 {
			return InStock
		}
		return OutOfStock
	}
	return OutOfStock
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
