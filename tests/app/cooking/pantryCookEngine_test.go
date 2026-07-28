package cooking

import (
	"errors"
	"testing"

	appCooking "github.com/tthung1997/buddy/app/cooking"
	coreCooking "github.com/tthung1997/buddy/core/cooking"
)

var _ coreCooking.ICookEngine = (*appCooking.PantryCookEngine)(nil)

func newEngine() *appCooking.PantryCookEngine {
	return appCooking.NewPantryCookEngine(appCooking.NewSimpleUnitConverter())
}

func pancakes() coreCooking.Recipe {
	return coreCooking.Recipe{
		Id:       "pancakes",
		Name:     "Pancakes",
		Servings: 2,
		Ingredients: []coreCooking.RecipeIngredient{
			{IngredientId: "flour", Amount: 200, Unit: coreCooking.Gram},
			{IngredientId: "milk", Amount: 300, Unit: coreCooking.Milliliter},
			{IngredientId: "egg", Amount: 2, Unit: coreCooking.Piece},
			{IngredientId: "maple-syrup", Amount: 2, Unit: coreCooking.Tablespoon, Optional: true},
		},
	}
}

func toast() coreCooking.Recipe {
	return coreCooking.Recipe{
		Id:       "toast",
		Name:     "Toast",
		Servings: 1,
		Ingredients: []coreCooking.RecipeIngredient{
			{IngredientId: "bread", Amount: 2, Unit: coreCooking.Piece},
		},
	}
}

func findMatch(t *testing.T, matches []coreCooking.RecipeMatch, id string) coreCooking.RecipeMatch {
	t.Helper()
	for _, match := range matches {
		if match.Recipe.Id == id {
			return match
		}
	}
	t.Fatalf("expected a match for %s, got %d matches", id, len(matches))
	return coreCooking.RecipeMatch{}
}

func TestPantryCookEngine_Match_CookableWhenPantryCovers(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 1, Unit: coreCooking.Kilogram},
		{IngredientId: "milk", Amount: 1, Unit: coreCooking.Liter},
		{IngredientId: "egg", Amount: 6, Unit: coreCooking.Piece},
	}

	matches := newEngine().Match([]coreCooking.Recipe{pancakes()}, pantry, coreCooking.MatchOptions{})
	match := findMatch(t, matches, "pancakes")

	if !match.Cookable {
		t.Fatalf("expected pancakes to be cookable, missing: %v", match.Missing)
	}
	assertClose(t, match.Coverage, 1, "coverage")
	if len(match.Have) != 3 {
		t.Errorf("expected 3 matched ingredients, got %d", len(match.Have))
	}
	if len(match.Missing) != 1 || match.Missing[0].IngredientId != "maple-syrup" {
		t.Errorf("expected only the optional syrup to be missing, got %v", match.Missing)
	}
	if !match.Missing[0].Optional {
		t.Error("expected the missing syrup to be flagged optional")
	}
}

func TestPantryCookEngine_Match_ReportsAbsentAndInsufficient(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 100, Unit: coreCooking.Gram},
		{IngredientId: "egg", Amount: 6, Unit: coreCooking.Piece},
	}

	matches := newEngine().Match([]coreCooking.Recipe{pancakes()}, pantry, coreCooking.MatchOptions{})
	match := findMatch(t, matches, "pancakes")

	if match.Cookable {
		t.Fatal("expected pancakes not to be cookable")
	}
	assertClose(t, match.Coverage, 1.0/3.0, "coverage")

	reasons := map[string]coreCooking.MissingReason{}
	amounts := map[string]float64{}
	for _, missing := range match.Missing {
		reasons[missing.IngredientId] = missing.Reason
		amounts[missing.IngredientId] = missing.Amount
	}

	if reasons["flour"] != coreCooking.MissingReasonInsufficient {
		t.Errorf("expected flour to be insufficient, got %s", reasons["flour"])
	}
	assertClose(t, amounts["flour"], 100, "missing flour")
	if reasons["milk"] != coreCooking.MissingReasonAbsent {
		t.Errorf("expected milk to be absent, got %s", reasons["milk"])
	}
	assertClose(t, amounts["milk"], 300, "missing milk")
}

func TestPantryCookEngine_Match_UnitMismatchDoesNotBreakTheMatch(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 1, Unit: coreCooking.Kilogram},
		{IngredientId: "milk", Amount: 1, Unit: coreCooking.Liter},
		{IngredientId: "egg", Amount: 500, Unit: coreCooking.Gram},
	}

	matches := newEngine().Match([]coreCooking.Recipe{pancakes()}, pantry, coreCooking.MatchOptions{})
	match := findMatch(t, matches, "pancakes")

	if match.Cookable {
		t.Fatal("expected an unmeasurable ingredient to block cooking")
	}
	found := false
	for _, missing := range match.Missing {
		if missing.IngredientId == "egg" {
			found = true
			if missing.Reason != coreCooking.MissingReasonUnitMismatch {
				t.Errorf("expected a unit mismatch for eggs, got %s", missing.Reason)
			}
		}
	}
	if !found {
		t.Error("expected eggs to be reported as missing")
	}
	assertClose(t, match.Coverage, 2.0/3.0, "coverage")
}

func TestPantryCookEngine_Match_IncludeOptional(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 1, Unit: coreCooking.Kilogram},
		{IngredientId: "milk", Amount: 1, Unit: coreCooking.Liter},
		{IngredientId: "egg", Amount: 6, Unit: coreCooking.Piece},
	}

	matches := newEngine().Match([]coreCooking.Recipe{pancakes()}, pantry, coreCooking.MatchOptions{IncludeOptional: true})
	match := findMatch(t, matches, "pancakes")

	if match.Cookable {
		t.Fatal("expected the optional ingredient to block cooking when included")
	}
	assertClose(t, match.Coverage, 3.0/4.0, "coverage")
}

func TestPantryCookEngine_Match_ScalesToRequestedServings(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 300, Unit: coreCooking.Gram},
		{IngredientId: "milk", Amount: 1, Unit: coreCooking.Liter},
		{IngredientId: "egg", Amount: 6, Unit: coreCooking.Piece},
	}

	matches := newEngine().Match([]coreCooking.Recipe{pancakes()}, pantry, coreCooking.MatchOptions{Servings: 4})
	match := findMatch(t, matches, "pancakes")

	if match.Servings != 4 {
		t.Errorf("expected 4 servings, got %d", match.Servings)
	}
	if match.Cookable {
		t.Fatal("expected doubling the recipe to exceed the available flour")
	}
	for _, missing := range match.Missing {
		if missing.IngredientId == "flour" {
			assertClose(t, missing.Amount, 100, "missing flour after scaling")
		}
	}
}

func TestPantryCookEngine_Match_SortsCookableFirstThenCoverage(t *testing.T) {
	empty := coreCooking.Recipe{
		Id:       "soup",
		Name:     "Soup",
		Servings: 1,
		Ingredients: []coreCooking.RecipeIngredient{
			{IngredientId: "stock", Amount: 1, Unit: coreCooking.Liter},
			{IngredientId: "carrot", Amount: 2, Unit: coreCooking.Piece},
		},
	}
	partial := coreCooking.Recipe{
		Id:       "omelette",
		Name:     "Omelette",
		Servings: 1,
		Ingredients: []coreCooking.RecipeIngredient{
			{IngredientId: "egg", Amount: 3, Unit: coreCooking.Piece},
			{IngredientId: "cheese", Amount: 50, Unit: coreCooking.Gram},
		},
	}
	pantry := []coreCooking.PantryItem{
		{IngredientId: "bread", Amount: 10, Unit: coreCooking.Piece},
		{IngredientId: "egg", Amount: 6, Unit: coreCooking.Piece},
	}

	matches := newEngine().Match([]coreCooking.Recipe{empty, partial, toast()}, pantry, coreCooking.MatchOptions{})

	if len(matches) != 3 {
		t.Fatalf("expected 3 matches, got %d", len(matches))
	}
	if matches[0].Recipe.Id != "toast" {
		t.Errorf("expected the cookable recipe first, got %s", matches[0].Recipe.Id)
	}
	if matches[1].Recipe.Id != "omelette" {
		t.Errorf("expected the better covered recipe second, got %s", matches[1].Recipe.Id)
	}
	if matches[2].Recipe.Id != "soup" {
		t.Errorf("expected the empty match last, got %s", matches[2].Recipe.Id)
	}
}

func TestPantryCookEngine_Match_Filters(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "bread", Amount: 10, Unit: coreCooking.Piece},
	}
	recipes := []coreCooking.Recipe{pancakes(), toast()}

	cookable := newEngine().Match(recipes, pantry, coreCooking.MatchOptions{OnlyCookable: true})
	if len(cookable) != 1 || cookable[0].Recipe.Id != "toast" {
		t.Fatalf("expected only toast, got %v", cookable)
	}

	covered := newEngine().Match(recipes, pantry, coreCooking.MatchOptions{MinCoverage: 0.5})
	if len(covered) != 1 || covered[0].Recipe.Id != "toast" {
		t.Fatalf("expected only toast above the coverage floor, got %v", covered)
	}
}

func TestPantryCookEngine_Match_MergesRepeatedIngredients(t *testing.T) {
	recipe := coreCooking.Recipe{
		Id:       "bread",
		Name:     "Bread",
		Servings: 1,
		Steps: []coreCooking.Step{
			{Order: 1, Ingredients: []coreCooking.RecipeIngredient{{IngredientId: "flour", Amount: 400, Unit: coreCooking.Gram}}},
			{Order: 2, Ingredients: []coreCooking.RecipeIngredient{{IngredientId: "flour", Amount: 0.1, Unit: coreCooking.Kilogram}}},
		},
	}
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 450, Unit: coreCooking.Gram},
	}

	match := findMatch(t, newEngine().Match([]coreCooking.Recipe{recipe}, pantry, coreCooking.MatchOptions{}), "bread")
	if match.Cookable {
		t.Fatal("expected the merged 500 g requirement to exceed 450 g of flour")
	}
	if len(match.Missing) != 1 {
		t.Fatalf("expected a single merged requirement, got %v", match.Missing)
	}
	assertClose(t, match.Missing[0].Amount, 50, "missing flour")
}

func TestPantryCookEngine_ShoppingList_AggregatesAcrossPlans(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 250, Unit: coreCooking.Gram},
		{IngredientId: "egg", Amount: 3, Unit: coreCooking.Piece},
	}
	plans := []coreCooking.RecipePlan{
		{Recipe: pancakes()},
		{Recipe: pancakes(), Servings: 4},
	}

	missing := newEngine().ShoppingList(plans, pantry)

	amounts := map[string]coreCooking.MissingIngredient{}
	for _, item := range missing {
		amounts[item.IngredientId] = item
	}

	// 200 g + 400 g needed against 250 g in the pantry.
	assertClose(t, amounts["flour"].Amount, 350, "missing flour")
	// 2 + 4 eggs needed against 3 in the pantry.
	assertClose(t, amounts["egg"].Amount, 3, "missing eggs")
	// 300 ml + 600 ml needed, none in the pantry.
	assertClose(t, amounts["milk"].Amount, 900, "missing milk")

	if _, exists := amounts["maple-syrup"]; exists {
		t.Error("expected optional ingredients to stay off the shopping list")
	}
}

func TestPantryCookEngine_ShoppingList_LeavesPantryUntouched(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 250, Unit: coreCooking.Gram},
	}

	newEngine().ShoppingList([]coreCooking.RecipePlan{{Recipe: pancakes()}}, pantry)

	assertClose(t, pantry[0].Amount, 250, "pantry after building a shopping list")
}

func TestPantryCookEngine_Consume(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 1, Unit: coreCooking.Kilogram},
		{IngredientId: "milk", Amount: 1, Unit: coreCooking.Liter},
		{IngredientId: "egg", Amount: 2, Unit: coreCooking.Piece},
	}

	updated, err := newEngine().Consume(pantry, pancakes(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	remaining := map[string]coreCooking.PantryItem{}
	for _, item := range updated {
		remaining[item.IngredientId] = item
	}

	assertClose(t, remaining["flour"].Amount, 0.8, "remaining flour")
	assertClose(t, remaining["milk"].Amount, 0.7, "remaining milk")
	if _, exists := remaining["egg"]; exists {
		t.Error("expected fully used ingredients to leave the pantry")
	}
	if remaining["flour"].UpdatedDateTime.IsZero() {
		t.Error("expected consumed items to be stamped")
	}
	assertClose(t, pantry[0].Amount, 1, "original pantry untouched")
}

func TestPantryCookEngine_Consume_Insufficient(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 100, Unit: coreCooking.Gram},
		{IngredientId: "milk", Amount: 1, Unit: coreCooking.Liter},
		{IngredientId: "egg", Amount: 2, Unit: coreCooking.Piece},
	}

	if _, err := newEngine().Consume(pantry, pancakes(), 2); !errors.Is(err, coreCooking.ErrInsufficientPantry) {
		t.Fatalf("expected ErrInsufficientPantry, got %v", err)
	}
}

func TestPantryCookEngine_Consume_MissingIngredient(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 1, Unit: coreCooking.Kilogram},
	}

	if _, err := newEngine().Consume(pantry, pancakes(), 2); !errors.Is(err, coreCooking.ErrInsufficientPantry) {
		t.Fatalf("expected ErrInsufficientPantry, got %v", err)
	}
}

func TestPantryCookEngine_Consume_IncompatibleUnit(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "bread", Amount: 500, Unit: coreCooking.Gram},
	}

	if _, err := newEngine().Consume(pantry, toast(), 1); !errors.Is(err, coreCooking.ErrIncompatibleUnits) {
		t.Fatalf("expected ErrIncompatibleUnits, got %v", err)
	}
}

func TestPantryCookEngine_Consume_ScalesServings(t *testing.T) {
	pantry := []coreCooking.PantryItem{
		{IngredientId: "bread", Amount: 10, Unit: coreCooking.Piece},
	}

	updated, err := newEngine().Consume(pantry, toast(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	assertClose(t, updated[0].Amount, 4, "remaining bread")
}
