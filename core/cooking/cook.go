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
)

type MatchedIngredient struct {
	IngredientId string  `json:"ingredientId"`
	Required     float64 `json:"required"`
	Available    float64 `json:"available"`
	Unit         Unit    `json:"unit"`
	Optional     bool    `json:"optional,omitempty"`
}

type MissingIngredient struct {
	IngredientId string        `json:"ingredientId"`
	Amount       float64       `json:"amount"`
	Unit         Unit          `json:"unit"`
	Optional     bool          `json:"optional,omitempty"`
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
	Match(recipes []Recipe, pantry []PantryItem, options MatchOptions) []RecipeMatch
	// ShoppingList aggregates everything the pantry cannot cover for the given
	// plans into one deduplicated list.
	ShoppingList(plans []RecipePlan, pantry []PantryItem) []MissingIngredient
	// Consume deducts the ingredients a recipe uses from the pantry and returns
	// the updated pantry.
	Consume(pantry []PantryItem, recipe Recipe, servings int) ([]PantryItem, error)
}
