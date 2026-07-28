package cooking

import (
	"errors"
	"testing"
	"time"

	"github.com/tthung1997/buddy/core/cooking"
)

func pancakes() cooking.Recipe {
	return cooking.Recipe{
		Id:       "pancakes",
		Name:     "Pancakes",
		Servings: 2,
		Ingredients: []cooking.RecipeIngredient{
			{IngredientId: "flour", Amount: 200, Unit: cooking.Gram},
			{IngredientId: "milk", Amount: 300, Unit: cooking.Milliliter},
			{IngredientId: "egg", Amount: 2, Unit: cooking.Piece},
			{IngredientId: "maple-syrup", Amount: 2, Unit: cooking.Tablespoon, Optional: true},
		},
		Steps: []cooking.Step{
			{Order: 2, Description: "Rest the batter", Duration: 30, DurationUnit: cooking.Minute},
			{Order: 1, Description: "Whisk everything", Duration: 5, DurationUnit: cooking.Minute,
				Ingredients: []cooking.RecipeIngredient{
					{IngredientId: "flour", Amount: 200, Unit: cooking.Gram},
				}},
			{Order: 3, Description: "Fry", Duration: 10, DurationUnit: cooking.Minute},
		},
	}
}

func TestRecipe_OrderedSteps(t *testing.T) {
	recipe := pancakes()
	steps := recipe.OrderedSteps()

	for i, step := range steps {
		if step.Order != i+1 {
			t.Fatalf("expected step %d at position %d, got step %d", i+1, i, step.Order)
		}
	}
	if recipe.Steps[0].Order != 2 {
		t.Error("expected OrderedSteps to leave the original recipe untouched")
	}
}

func TestRecipe_TotalDuration(t *testing.T) {
	total, err := pancakes().TotalDuration()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 45*time.Minute {
		t.Errorf("expected 45m, got %v", total)
	}
}

func TestRecipe_TotalDuration_IgnoresStepsWithoutDuration(t *testing.T) {
	recipe := cooking.Recipe{
		Servings: 1,
		Steps: []cooking.Step{
			{Order: 1, Description: "Season to taste"},
			{Order: 2, Description: "Bake", Duration: 20, DurationUnit: cooking.Minute},
		},
	}

	total, err := recipe.TotalDuration()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if total != 20*time.Minute {
		t.Errorf("expected 20m, got %v", total)
	}
}

func TestRecipe_TotalDuration_UnknownDurationUnit(t *testing.T) {
	recipe := cooking.Recipe{
		Servings: 1,
		Steps: []cooking.Step{
			{Order: 1, Description: "Wait", Duration: 3, DurationUnit: cooking.DurationUnit("fortnight")},
		},
	}

	if _, err := recipe.TotalDuration(); !errors.Is(err, cooking.ErrUnknownDurationUnit) {
		t.Fatalf("expected ErrUnknownDurationUnit, got %v", err)
	}
}

func TestRecipe_ScaleTo(t *testing.T) {
	scaled, err := pancakes().ScaleTo(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if scaled.Servings != 3 {
		t.Errorf("expected 3 servings, got %d", scaled.Servings)
	}
	assertClose(t, scaled.Ingredients[0].Amount, 300, "scaled flour")
	assertClose(t, scaled.Ingredients[1].Amount, 450, "scaled milk")
	assertClose(t, scaled.Ingredients[2].Amount, 3, "scaled eggs")

	for _, step := range scaled.Steps {
		for _, ingredient := range step.Ingredients {
			if ingredient.IngredientId == "flour" {
				assertClose(t, ingredient.Amount, 300, "scaled step flour")
			}
		}
	}

	original := pancakes()
	assertClose(t, original.Ingredients[0].Amount, 200, "original flour untouched")
}

func TestRecipe_ScaleTo_InvalidServings(t *testing.T) {
	if _, err := pancakes().ScaleTo(0); !errors.Is(err, cooking.ErrInvalidServings) {
		t.Fatalf("expected ErrInvalidServings, got %v", err)
	}

	unservable := cooking.Recipe{Name: "Mystery", Servings: 0}
	if _, err := unservable.ScaleTo(2); !errors.Is(err, cooking.ErrInvalidServings) {
		t.Fatalf("expected ErrInvalidServings, got %v", err)
	}
}

func TestRecipe_RequiredIngredients_FallsBackToSteps(t *testing.T) {
	recipe := cooking.Recipe{
		Servings: 1,
		Steps: []cooking.Step{
			{Order: 2, Ingredients: []cooking.RecipeIngredient{{IngredientId: "salt", Amount: 1, Unit: cooking.Teaspoon}}},
			{Order: 1, Ingredients: []cooking.RecipeIngredient{{IngredientId: "rice", Amount: 1, Unit: cooking.Cup}}},
		},
	}

	required := recipe.RequiredIngredients()
	if len(required) != 2 {
		t.Fatalf("expected 2 required ingredients, got %d", len(required))
	}
	if required[0].IngredientId != "rice" || required[1].IngredientId != "salt" {
		t.Errorf("expected step order to be respected, got %v", required)
	}
}

func TestRecipe_RequiredIngredients_PrefersRecipeLevelList(t *testing.T) {
	required := pancakes().RequiredIngredients()
	if len(required) != 4 {
		t.Fatalf("expected 4 required ingredients, got %d", len(required))
	}
	for _, ingredient := range required {
		if ingredient.IngredientId == "flour" && ingredient.Amount != 200 {
			t.Errorf("expected the recipe level amount to win, got %v", ingredient.Amount)
		}
	}
}

func TestRecipe_RequiredIngredients_AddsStepOnlyIngredients(t *testing.T) {
	recipe := pancakes()
	recipe.Steps = append(recipe.Steps, cooking.Step{
		Order:       4,
		Description: "Fry in butter",
		Ingredients: []cooking.RecipeIngredient{
			{IngredientId: "butter", Amount: 20, Unit: cooking.Gram},
		},
	})

	required := recipe.RequiredIngredients()
	if len(required) != 5 {
		t.Fatalf("expected the step only ingredient to be included, got %d", len(required))
	}
	if required[4].IngredientId != "butter" {
		t.Errorf("expected butter to be appended, got %s", required[4].IngredientId)
	}
}

func TestIngredient_Allows(t *testing.T) {
	flour := cooking.Ingredient{Id: "flour", Name: "Flour", DefaultUnit: cooking.Gram}
	if !flour.Allows(cooking.Kilogram) {
		t.Error("expected a mass ingredient to allow kilograms")
	}
	if flour.Allows(cooking.Liter) {
		t.Error("expected a mass ingredient to reject liters")
	}
	if flour.Allows(cooking.Unit("furlong")) {
		t.Error("expected an unknown unit to be rejected")
	}

	egg := cooking.Ingredient{
		Id:           "egg",
		Name:         "Egg",
		DefaultUnit:  cooking.Piece,
		AllowedUnits: []cooking.Unit{cooking.Piece, cooking.Gram},
	}
	if !egg.Allows(cooking.Gram) {
		t.Error("expected an explicit allow list to be honoured")
	}
	if egg.Allows(cooking.Cup) {
		t.Error("expected units outside the allow list to be rejected")
	}
}
