package cooking

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/tthung1997/buddy/app/cooking/repository"
	coreCooking "github.com/tthung1997/buddy/core/cooking"
)

var (
	_ coreCooking.IRecipeRepository     = (*repository.LocalRecipeRepository)(nil)
	_ coreCooking.IIngredientRepository = (*repository.LocalIngredientRepository)(nil)
	_ coreCooking.IPantryRepository     = (*repository.LocalPantryRepository)(nil)
)

func TestLocalRecipeRepository_SeedsAndRoundTrips(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), ".db", "recipes.json")

	repo, err := repository.NewLocalRecipeRepository(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(filePath); err != nil {
		t.Fatalf("expected the repository to seed %s: %v", filePath, err)
	}

	recipes, err := repo.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recipes) != 0 {
		t.Fatalf("expected an empty repository, got %d recipes", len(recipes))
	}

	recipe := coreCooking.Recipe{
		Id:       "toast",
		Name:     "Toast",
		Servings: 1,
		Ingredients: []coreCooking.RecipeIngredient{
			{IngredientId: "bread", Amount: 2, Unit: coreCooking.Piece},
		},
		Steps: []coreCooking.Step{
			{Order: 1, Description: "Toast the bread", Duration: 3, DurationUnit: coreCooking.Minute},
		},
		CreatedDateTime: time.Now().UTC().Truncate(time.Second),
	}
	if err := repo.CreateOrUpdate(recipe); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stored, err := repo.Get("toast")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stored.Name != "Toast" || len(stored.Steps) != 1 || stored.Steps[0].DurationUnit != coreCooking.Minute {
		t.Errorf("unexpected round trip result: %+v", stored)
	}

	// A second repository over the same file must observe the write.
	reopened, err := repository.NewLocalRecipeRepository(filePath)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	recipes, err = reopened.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recipes) != 1 || recipes[0].Id != "toast" {
		t.Fatalf("expected the stored recipe to persist, got %v", recipes)
	}

	if err := repo.Delete("toast"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := repo.Get("toast"); !errors.Is(err, coreCooking.ErrRecipeNotFound) {
		t.Fatalf("expected ErrRecipeNotFound, got %v", err)
	}
}

func TestLocalRecipeRepository_ListIsSortedByName(t *testing.T) {
	repo, err := repository.NewLocalRecipeRepository(filepath.Join(t.TempDir(), "recipes.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, recipe := range []coreCooking.Recipe{
		{Id: "c", Name: "Cassoulet", Servings: 4},
		{Id: "a", Name: "Arepas", Servings: 2},
		{Id: "b", Name: "Borscht", Servings: 6},
	} {
		if err := repo.CreateOrUpdate(recipe); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	recipes, err := repo.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"Arepas", "Borscht", "Cassoulet"}
	for i, name := range expected {
		if recipes[i].Name != name {
			t.Fatalf("expected %v, got %v at %d", name, recipes[i].Name, i)
		}
	}
}

func TestLocalRecipeRepository_RejectsEmptyId(t *testing.T) {
	repo, err := repository.NewLocalRecipeRepository(filepath.Join(t.TempDir(), "recipes.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.CreateOrUpdate(coreCooking.Recipe{Name: "Nameless"}); err == nil {
		t.Fatal("expected an error when storing a recipe without an id")
	}
}

func TestLocalIngredientRepository_RoundTrip(t *testing.T) {
	repo, err := repository.NewLocalIngredientRepository(filepath.Join(t.TempDir(), "ingredients.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	ingredient := coreCooking.Ingredient{
		Id:           "flour",
		Name:         "Flour",
		Category:     "Baking",
		DefaultUnit:  coreCooking.Gram,
		AllowedUnits: []coreCooking.Unit{coreCooking.Gram, coreCooking.Kilogram},
	}
	if err := repo.CreateOrUpdate(ingredient); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	stored, err := repo.Get("flour")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stored.DefaultUnit != coreCooking.Gram || len(stored.AllowedUnits) != 2 {
		t.Errorf("unexpected round trip result: %+v", stored)
	}

	if _, err := repo.Get("sugar"); !errors.Is(err, coreCooking.ErrIngredientNotFound) {
		t.Fatalf("expected ErrIngredientNotFound, got %v", err)
	}
}

func TestLocalPantryRepository_ReplaceAll(t *testing.T) {
	repo, err := repository.NewLocalPantryRepository(filepath.Join(t.TempDir(), "pantry.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.CreateOrUpdate(coreCooking.PantryItem{IngredientId: "milk", Amount: 1, Unit: coreCooking.Liter}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.ReplaceAll([]coreCooking.PantryItem{
		{IngredientId: "flour", Amount: 500, Unit: coreCooking.Gram},
		{IngredientId: "egg", Amount: 6, Unit: coreCooking.Piece},
	}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	items, err := repo.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("expected the pantry to be replaced, got %d items", len(items))
	}
	if _, err := repo.Get("milk"); !errors.Is(err, coreCooking.ErrPantryItemNotFound) {
		t.Fatalf("expected the previous entry to be gone, got %v", err)
	}

	if err := repo.ReplaceAll(nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	items, err = repo.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected an empty pantry, got %d items", len(items))
	}
}

func TestLocalIngredientRepository_ReplaceAll_RejectsEmptyId(t *testing.T) {
	repo, err := repository.NewLocalIngredientRepository(filepath.Join(t.TempDir(), "ingredients.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.CreateOrUpdate(coreCooking.Ingredient{Id: "flour", Name: "Flour", DefaultUnit: coreCooking.Gram}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.ReplaceAll([]coreCooking.Ingredient{{Name: "Nameless"}}); err == nil {
		t.Fatal("expected an error when replacing with an ingredient without an id")
	}

	ingredients, err := repo.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(ingredients) != 1 || ingredients[0].Id != "flour" {
		t.Fatalf("expected the rejected write to leave the catalog untouched, got %v", ingredients)
	}
}

func TestLocalPantryRepository_RoundTrip(t *testing.T) {
	repo, err := repository.NewLocalPantryRepository(filepath.Join(t.TempDir(), "pantry.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := repo.CreateOrUpdate(coreCooking.PantryItem{IngredientId: "milk", Amount: 1, Unit: coreCooking.Liter}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := repo.CreateOrUpdate(coreCooking.PantryItem{IngredientId: "milk", Amount: 2, Unit: coreCooking.Liter}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	items, err := repo.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected the second write to replace the first, got %d items", len(items))
	}
	assertClose(t, items[0].Amount, 2, "stored amount")

	if _, err := repo.Get("eggs"); !errors.Is(err, coreCooking.ErrPantryItemNotFound) {
		t.Fatalf("expected ErrPantryItemNotFound, got %v", err)
	}

	if err := repo.Delete("milk"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	items, err = repo.List()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Fatalf("expected an empty pantry, got %d items", len(items))
	}
}
