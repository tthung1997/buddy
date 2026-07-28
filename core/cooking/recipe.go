package cooking

import (
	"errors"
	"fmt"
	"sort"
	"time"
)

var (
	ErrRecipeNotFound  = errors.New("recipe not found")
	ErrInvalidServings = errors.New("servings must be greater than zero")
)

type RecipeIngredient struct {
	IngredientId string  `json:"ingredientId"`
	Amount       float64 `json:"amount"`
	Unit         Unit    `json:"unit"`
	Optional     bool    `json:"optional,omitempty"`
	Note         string  `json:"note,omitempty"`
}

type Step struct {
	Order        int                `json:"order"`
	Description  string             `json:"description"`
	ImageURL     string             `json:"imageUrl,omitempty"`
	Duration     float64            `json:"duration,omitempty"`
	DurationUnit DurationUnit       `json:"durationUnit,omitempty"`
	Ingredients  []RecipeIngredient `json:"ingredients,omitempty"`
}

// TotalDuration returns how long the step takes. A step without a duration unit
// is treated as instantaneous rather than as an error.
func (s Step) TotalDuration() (time.Duration, error) {
	if s.Duration == 0 || s.DurationUnit == "" {
		return 0, nil
	}
	return s.DurationUnit.ToDuration(s.Duration)
}

type Recipe struct {
	Id          string             `json:"id"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	ImageURL    string             `json:"imageUrl,omitempty"`
	Servings    int                `json:"servings"`
	Tags        []string           `json:"tags,omitempty"`
	Ingredients []RecipeIngredient `json:"ingredients,omitempty"`
	Steps       []Step             `json:"steps,omitempty"`
	Rating      int                `json:"rating,omitempty"`
	Notes       string             `json:"notes,omitempty"`

	CreatedDateTime    time.Time  `json:"createdDateTime"`
	UpdatedDateTime    time.Time  `json:"updatedDateTime"`
	LastCookedDateTime *time.Time `json:"lastCookedDateTime,omitempty"`
}

// OrderedSteps returns a copy of the steps sorted by their declared order.
func (r Recipe) OrderedSteps() []Step {
	steps := make([]Step, len(r.Steps))
	copy(steps, r.Steps)
	sort.SliceStable(steps, func(i, j int) bool {
		return steps[i].Order < steps[j].Order
	})
	return steps
}

// TotalDuration sums the duration of every step in the recipe.
func (r Recipe) TotalDuration() (time.Duration, error) {
	var total time.Duration
	for _, step := range r.Steps {
		duration, err := step.TotalDuration()
		if err != nil {
			return 0, fmt.Errorf("step %d: %w", step.Order, err)
		}
		total += duration
	}
	return total, nil
}

// RequiredIngredients returns everything the recipe consumes. The recipe level
// list is authoritative for the amounts, and any ingredient that only appears
// inside a step is added on top of it so nothing goes unnoticed.
func (r Recipe) RequiredIngredients() []RecipeIngredient {
	ingredients := make([]RecipeIngredient, 0, len(r.Ingredients))
	declared := make(map[string]bool, len(r.Ingredients))
	for _, ingredient := range r.Ingredients {
		ingredients = append(ingredients, ingredient)
		declared[ingredient.IngredientId] = true
	}

	for _, step := range r.OrderedSteps() {
		for _, ingredient := range step.Ingredients {
			if declared[ingredient.IngredientId] {
				continue
			}
			ingredients = append(ingredients, ingredient)
		}
	}

	return ingredients
}

// ScaleTo returns a copy of the recipe with every ingredient amount scaled to
// the requested number of servings.
func (r Recipe) ScaleTo(servings int) (Recipe, error) {
	if servings <= 0 {
		return Recipe{}, fmt.Errorf("%w: requested %d", ErrInvalidServings, servings)
	}
	if r.Servings <= 0 {
		return Recipe{}, fmt.Errorf("%w: recipe %q declares %d", ErrInvalidServings, r.Name, r.Servings)
	}

	factor := float64(servings) / float64(r.Servings)
	scaled := r
	scaled.Servings = servings
	scaled.Tags = append([]string(nil), r.Tags...)
	scaled.Ingredients = scaleIngredients(r.Ingredients, factor)

	scaled.Steps = make([]Step, len(r.Steps))
	for i, step := range r.Steps {
		step.Ingredients = scaleIngredients(step.Ingredients, factor)
		scaled.Steps[i] = step
	}

	return scaled, nil
}

func scaleIngredients(ingredients []RecipeIngredient, factor float64) []RecipeIngredient {
	if ingredients == nil {
		return nil
	}
	scaled := make([]RecipeIngredient, len(ingredients))
	for i, ingredient := range ingredients {
		ingredient.Amount *= factor
		scaled[i] = ingredient
	}
	return scaled
}

type IRecipeRepository interface {
	List() ([]Recipe, error)
	Get(id string) (Recipe, error)
	CreateOrUpdate(Recipe) error
	Delete(id string) error
}
