package cooking

import (
	"github.com/tthung1997/buddy/core/cooking"
)

// IngredientLine is a recipe ingredient with its catalog name resolved and its
// amount already formatted for display.
type IngredientLine struct {
	IngredientId string
	Name         string
	Amount       float64
	AmountLabel  string
	Unit         cooking.Unit
	Optional     bool
	Note         string
	ImageURL     string
}

// MissingLine explains why an ingredient keeps a recipe off the menu.
type MissingLine struct {
	IngredientLine
	Reason      cooking.MissingReason
	ReasonLabel string
}

type StepView struct {
	Order         int
	Description   string
	ImageURL      string
	DurationLabel string
	Ingredients   []IngredientLine
}

type RecipeView struct {
	Id                 string
	Name               string
	Description        string
	ImageURL           string
	Servings           int
	Tags               []string
	Rating             int
	Notes              string
	Ingredients        []IngredientLine
	Steps              []StepView
	TotalDurationLabel string
	LastCookedLabel    string
	StepCount          int
	IngredientCount    int
}

type MatchView struct {
	Recipe          RecipeView
	Servings        int
	Coverage        float64
	CoveragePercent int
	Cookable        bool
	Have            []IngredientLine
	Missing         []MissingLine
	MissingCount    int
	StatusLabel     string
	StatusClass     string
}

type IndexPageData struct {
	Error           string
	Success         string
	RecipeCount     int
	IngredientCount int
	PantryCount     int
	CookableCount   int
	Highlights      []MatchView
	NextUp          []MissingLine
}

type RecipesPageData struct {
	Error   string
	Success string
	Query   string
	Tag     string
	Tags    []string
	Recipes []RecipeView
	Count   int
	Total   int
}

type RecipePageData struct {
	Error    string
	Success  string
	Recipe   RecipeView
	Servings int
	Match    MatchView
	HasMatch bool
}

type CookPageData struct {
	Error           string
	Success         string
	Servings        int
	IncludeOptional bool
	OnlyCookable    bool
	Matches         []MatchView
	CookableCount   int
	TotalCount      int
}

type PantryPageData struct {
	Error         string
	Success       string
	Units         []cooking.Unit
	IngredientIds []string
}

type EditPageData struct {
	Error         string
	IsNew         bool
	RecipeId      string
	RecipeName    string
	RecipeJSON    string
	CatalogJSON   string
	UnitsJSON     string
	DurationsJSON string
}
