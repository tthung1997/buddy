package cooking

import "errors"

var ErrInsufficientPantry = errors.New("not enough of an ingredient in the pantry")

type MissingReason string

const (
	// MissingReasonAbsent means the pantry holds none of the ingredient.
	MissingReasonAbsent MissingReason = "ABSENT"
	// MissingReasonInsufficient means the pantry holds some, but not enough.
	MissingReasonInsufficient MissingReason = "INSUFFICIENT"
	// MissingReasonUnitMismatch means the pantry holds the ingredient in a unit
	// that cannot be compared with the one the recipe asks for.
	MissingReasonUnitMismatch MissingReason = "UNIT_MISMATCH"
	// MissingReasonOutOfStock means a staple has run out. Staples carry no
	// amount, so this replaces absent and insufficient for them.
	MissingReasonOutOfStock MissingReason = "OUT_OF_STOCK"
	// MissingReasonRunningLow means a staple still cooks tonight but should be
	// restocked. It never blocks a recipe, so it only reaches shopping lists.
	MissingReasonRunningLow MissingReason = "RUNNING_LOW"
)

type MatchedIngredient struct {
	IngredientId string  `json:"ingredientId"`
	Required     float64 `json:"required"`
	Available    float64 `json:"available"`
	Unit         Unit    `json:"unit"`
	Optional     bool    `json:"optional,omitempty"`
	// Staple marks an ingredient satisfied by presence rather than by amount.
	Staple bool `json:"staple,omitempty"`
	// Stock carries the staple's level so a view can say "running low" instead
	// of printing an amount that was never tracked.
	Stock StockLevel `json:"stock,omitempty"`
}

type MissingIngredient struct {
	IngredientId string        `json:"ingredientId"`
	Amount       float64       `json:"amount,omitempty"`
	Unit         Unit          `json:"unit,omitempty"`
	Optional     bool          `json:"optional,omitempty"`
	Staple       bool          `json:"staple,omitempty"`
	Reason       MissingReason `json:"reason"`
}

type RecipeMatch struct {
	Recipe   Recipe              `json:"recipe"`
	Servings int                 `json:"servings"`
	Have     []MatchedIngredient `json:"have"`
	Missing  []MissingIngredient `json:"missing"`
	// Coverage is the share of required ingredients the pantry satisfies, from
	// 0 to 1. Recipes with no required ingredients score 1.
	Coverage float64 `json:"coverage"`
	Cookable bool    `json:"cookable"`
}

type MatchOptions struct {
	// Servings scales every recipe before matching. Zero keeps each recipe at
	// its own declared serving count.
	Servings int `json:"servings,omitempty"`
	// IncludeOptional treats optional ingredients as required.
	IncludeOptional bool `json:"includeOptional,omitempty"`
	// OnlyCookable drops recipes that cannot be cooked right now.
	OnlyCookable bool `json:"onlyCookable,omitempty"`
	// MinCoverage drops recipes below the given coverage share.
	MinCoverage float64 `json:"minCoverage,omitempty"`
}

type RecipePlan struct {
	Recipe   Recipe `json:"recipe"`
	Servings int    `json:"servings,omitempty"`
}

type ICookEngine interface {
	// Match joins the pantry against each recipe and reports what is available,
	// what is missing, and whether the recipe can be cooked right now. Results
	// are ordered with cookable recipes first, then by descending coverage.
	Match(recipes []Recipe, pantry Pantry, options MatchOptions) []RecipeMatch
	// ShoppingList aggregates everything the pantry cannot cover for the given
	// plans into one deduplicated list, including staples that have run out or
	// are running low.
	ShoppingList(plans []RecipePlan, pantry Pantry) []MissingIngredient
	// Consume deducts the ingredients a recipe uses from the pantry and returns
	// the updated items. Staples are left untouched.
	Consume(pantry Pantry, recipe Recipe, servings int) ([]PantryItem, error)
}
