package repository

import (
	"fmt"
	"sort"

	"github.com/tthung1997/buddy/core/cooking"
)

type LocalIngredientRepository struct {
	store *localStore[cooking.Ingredient]
}

func NewLocalIngredientRepository(filePath string) (*LocalIngredientRepository, error) {
	store, err := newLocalStore[cooking.Ingredient](filePath)
	if err != nil {
		return nil, err
	}
	return &LocalIngredientRepository{store: store}, nil
}

func (r *LocalIngredientRepository) List() ([]cooking.Ingredient, error) {
	items, err := r.store.all()
	if err != nil {
		return nil, err
	}

	ingredients := make([]cooking.Ingredient, 0, len(items))
	for _, ingredient := range items {
		ingredients = append(ingredients, ingredient)
	}
	sort.SliceStable(ingredients, func(i, j int) bool {
		if ingredients[i].Name != ingredients[j].Name {
			return ingredients[i].Name < ingredients[j].Name
		}
		return ingredients[i].Id < ingredients[j].Id
	})

	return ingredients, nil
}

func (r *LocalIngredientRepository) Get(id string) (cooking.Ingredient, error) {
	ingredient, exists, err := r.store.get(id)
	if err != nil {
		return cooking.Ingredient{}, err
	}
	if !exists {
		return cooking.Ingredient{}, fmt.Errorf("%w: %s", cooking.ErrIngredientNotFound, id)
	}
	return ingredient, nil
}

func (r *LocalIngredientRepository) CreateOrUpdate(ingredient cooking.Ingredient) error {
	return r.store.put(ingredient.Id, ingredient)
}

func (r *LocalIngredientRepository) Delete(id string) error {
	return r.store.remove(id)
}
