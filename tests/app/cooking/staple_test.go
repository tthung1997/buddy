package cooking

import (
	"errors"
	"testing"

	coreCooking "github.com/tthung1997/buddy/core/cooking"
)

// saltedSoup needs a pinch of a staple and a measured amount of everything else.
func saltedSoup() coreCooking.Recipe {
	return coreCooking.Recipe{
		Id:       "salted-soup",
		Name:     "Salted Soup",
		Servings: 2,
		Ingredients: []coreCooking.RecipeIngredient{
			{IngredientId: "carrot", Amount: 2, Unit: coreCooking.Piece},
			{IngredientId: "salt", Amount: 5, Unit: coreCooking.Gram},
		},
	}
}

func stapleKitchen(stock coreCooking.StockLevel) coreCooking.Pantry {
	items := []coreCooking.PantryItem{
		{IngredientId: "carrot", Amount: 6, Unit: coreCooking.Piece},
		{IngredientId: "salt", Stock: stock},
	}
	catalog := map[string]coreCooking.Ingredient{
		"carrot": {Id: "carrot", Name: "Carrot", DefaultUnit: coreCooking.Piece},
		"salt":   {Id: "salt", Name: "Salt", Staple: true},
	}
	return coreCooking.NewPantry(items, catalog)
}

func TestStaple_SatisfiesRegardlessOfAmount(t *testing.T) {
	// The recipe asks for 5 g of salt and the pantry records no amount at all.
	match := newEngine().Match([]coreCooking.Recipe{saltedSoup()}, stapleKitchen(coreCooking.InStock), coreCooking.MatchOptions{})[0]

	if !match.Cookable {
		t.Fatalf("expected an in stock staple to satisfy the recipe, missing: %v", match.Missing)
	}
	assertClose(t, match.Coverage, 1, "coverage")

	found := false
	for _, have := range match.Have {
		if have.IngredientId == "salt" {
			found = true
			if !have.Staple {
				t.Error("expected the salt line to be flagged as a staple")
			}
			if have.Stock != coreCooking.InStock {
				t.Errorf("expected the stock level to be carried through, got %q", have.Stock)
			}
		}
	}
	if !found {
		t.Error("expected salt to appear in the matched ingredients")
	}
}

func TestStaple_RunningLowStillCooks(t *testing.T) {
	match := newEngine().Match([]coreCooking.Recipe{saltedSoup()}, stapleKitchen(coreCooking.LowStock), coreCooking.MatchOptions{})[0]

	if !match.Cookable {
		t.Fatalf("expected a running low staple to still cook, missing: %v", match.Missing)
	}
	if len(match.Missing) != 0 {
		t.Errorf("expected nothing to be missing, got %v", match.Missing)
	}
}

func TestStaple_OutOfStockBlocksTheRecipe(t *testing.T) {
	match := newEngine().Match([]coreCooking.Recipe{saltedSoup()}, stapleKitchen(coreCooking.OutOfStock), coreCooking.MatchOptions{})[0]

	if match.Cookable {
		t.Fatal("expected an out of stock staple to block the recipe")
	}
	assertClose(t, match.Coverage, 0.5, "coverage")

	if len(match.Missing) != 1 {
		t.Fatalf("expected one missing ingredient, got %v", match.Missing)
	}
	if match.Missing[0].IngredientId != "salt" {
		t.Fatalf("expected salt to be missing, got %s", match.Missing[0].IngredientId)
	}
	if match.Missing[0].Reason != coreCooking.MissingReasonOutOfStock {
		t.Errorf("expected an out of stock reason, got %s", match.Missing[0].Reason)
	}
	if match.Missing[0].Amount != 0 {
		t.Errorf("expected a staple to carry no amount, got %v", match.Missing[0].Amount)
	}
}

func TestStaple_MissingPantryRowCountsAsOut(t *testing.T) {
	pantry := coreCooking.NewPantry(
		[]coreCooking.PantryItem{{IngredientId: "carrot", Amount: 6, Unit: coreCooking.Piece}},
		map[string]coreCooking.Ingredient{
			"carrot": {Id: "carrot", Name: "Carrot"},
			"salt":   {Id: "salt", Name: "Salt", Staple: true},
		},
	)

	match := newEngine().Match([]coreCooking.Recipe{saltedSoup()}, pantry, coreCooking.MatchOptions{})[0]
	if match.Cookable {
		t.Fatal("expected a staple with no pantry row to block the recipe")
	}
	if match.Missing[0].Reason != coreCooking.MissingReasonOutOfStock {
		t.Errorf("expected an out of stock reason, got %s", match.Missing[0].Reason)
	}
}

func TestStaple_NeverReportsAUnitMismatch(t *testing.T) {
	// The recipe measures salt by volume while the pantry tracks it by presence.
	recipe := saltedSoup()
	recipe.Ingredients[1].Unit = coreCooking.Teaspoon

	match := newEngine().Match([]coreCooking.Recipe{recipe}, stapleKitchen(coreCooking.InStock), coreCooking.MatchOptions{})[0]
	if !match.Cookable {
		t.Fatalf("expected presence to satisfy the staple whatever the unit, missing: %v", match.Missing)
	}
}

func TestStaple_ScalingDoesNotAffectIt(t *testing.T) {
	match := newEngine().Match([]coreCooking.Recipe{saltedSoup()}, stapleKitchen(coreCooking.InStock), coreCooking.MatchOptions{Servings: 6})[0]

	if !match.Cookable {
		t.Fatalf("expected the staple to survive scaling, missing: %v", match.Missing)
	}
}

func TestStaple_IsNotDeductedByCooking(t *testing.T) {
	pantry := stapleKitchen(coreCooking.InStock)

	updated, err := newEngine().Consume(pantry, saltedSoup(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	remaining := map[string]coreCooking.PantryItem{}
	for _, item := range updated {
		remaining[item.IngredientId] = item
	}

	salt, exists := remaining["salt"]
	if !exists {
		t.Fatal("expected the staple to survive cooking")
	}
	if salt.Stock != coreCooking.InStock {
		t.Errorf("expected the stock level to be untouched, got %q", salt.Stock)
	}
	if salt.Amount != 0 {
		t.Errorf("expected a staple to carry no amount, got %v", salt.Amount)
	}
	assertClose(t, remaining["carrot"].Amount, 4, "remaining carrots")
}

func TestStaple_OutOfStockBlocksCooking(t *testing.T) {
	if _, err := newEngine().Consume(stapleKitchen(coreCooking.OutOfStock), saltedSoup(), 2); !errors.Is(err, coreCooking.ErrInsufficientPantry) {
		t.Fatalf("expected ErrInsufficientPantry, got %v", err)
	}
}

func TestStaple_ShoppingListSkipsWhatIsInStock(t *testing.T) {
	missing := newEngine().ShoppingList(
		[]coreCooking.RecipePlan{{Recipe: saltedSoup()}},
		stapleKitchen(coreCooking.InStock),
	)

	for _, item := range missing {
		if item.IngredientId == "salt" {
			t.Errorf("expected an in stock staple to stay off the list, got %v", item)
		}
	}
}

func TestStaple_ShoppingListIncludesRunningLow(t *testing.T) {
	missing := newEngine().ShoppingList(
		[]coreCooking.RecipePlan{{Recipe: saltedSoup()}},
		stapleKitchen(coreCooking.LowStock),
	)

	found := false
	for _, item := range missing {
		if item.IngredientId != "salt" {
			continue
		}
		found = true
		if item.Reason != coreCooking.MissingReasonRunningLow {
			t.Errorf("expected a running low reason, got %s", item.Reason)
		}
		if item.Amount != 0 || item.Unit != "" {
			t.Errorf("expected a staple to carry no quantity, got %v %s", item.Amount, item.Unit)
		}
		if !item.Staple {
			t.Error("expected the line to be flagged as a staple")
		}
	}
	if !found {
		t.Error("expected a running low staple to reach the shopping list")
	}
}

func TestStaple_ShoppingListIncludesOutOfStock(t *testing.T) {
	missing := newEngine().ShoppingList(
		[]coreCooking.RecipePlan{{Recipe: saltedSoup()}},
		stapleKitchen(coreCooking.OutOfStock),
	)

	found := false
	for _, item := range missing {
		if item.IngredientId == "salt" {
			found = true
			if item.Reason != coreCooking.MissingReasonOutOfStock {
				t.Errorf("expected an out of stock reason, got %s", item.Reason)
			}
		}
	}
	if !found {
		t.Error("expected an out of stock staple to reach the shopping list")
	}
}

func TestStaple_ShoppingListRecordsAStapleOnce(t *testing.T) {
	// Two recipes both needing the same out of stock staple should buy it once.
	plans := []coreCooking.RecipePlan{
		{Recipe: saltedSoup()},
		{Recipe: saltedSoup(), Servings: 4},
	}

	missing := newEngine().ShoppingList(plans, stapleKitchen(coreCooking.OutOfStock))

	count := 0
	for _, item := range missing {
		if item.IngredientId == "salt" {
			count++
		}
	}
	if count != 1 {
		t.Fatalf("expected salt to be listed once, got %d", count)
	}
}

func TestStaple_OutOfStockWinsOverRunningLow(t *testing.T) {
	pantry := coreCooking.NewPantry(
		[]coreCooking.PantryItem{
			{IngredientId: "carrot", Amount: 20, Unit: coreCooking.Piece},
			{IngredientId: "salt", Stock: coreCooking.OutOfStock},
		},
		map[string]coreCooking.Ingredient{
			"carrot": {Id: "carrot", Name: "Carrot"},
			"salt":   {Id: "salt", Name: "Salt", Staple: true},
		},
	)

	missing := newEngine().ShoppingList([]coreCooking.RecipePlan{{Recipe: saltedSoup()}}, pantry)
	for _, item := range missing {
		if item.IngredientId == "salt" && item.Reason != coreCooking.MissingReasonOutOfStock {
			t.Errorf("expected the more urgent reason to win, got %s", item.Reason)
		}
	}
}

func TestStaple_LegacyAmountCountsAsInStock(t *testing.T) {
	// A pantry row written before the ingredient became a staple carries an
	// amount rather than a level; it should still read as in stock.
	pantry := coreCooking.NewPantry(
		[]coreCooking.PantryItem{
			{IngredientId: "carrot", Amount: 6, Unit: coreCooking.Piece},
			{IngredientId: "salt", Amount: 300, Unit: coreCooking.Gram},
		},
		map[string]coreCooking.Ingredient{
			"carrot": {Id: "carrot", Name: "Carrot"},
			"salt":   {Id: "salt", Name: "Salt", Staple: true},
		},
	)

	match := newEngine().Match([]coreCooking.Recipe{saltedSoup()}, pantry, coreCooking.MatchOptions{})[0]
	if !match.Cookable {
		t.Fatalf("expected a legacy amount to read as in stock, missing: %v", match.Missing)
	}
}
