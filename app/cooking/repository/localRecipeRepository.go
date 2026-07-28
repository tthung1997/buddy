package repository

import (
	"fmt"
	"sort"

	"github.com/tthung1997/buddy/core/cooking"
)

type LocalRecipeRepository struct {
	store *localStore[cooking.Recipe]
}

func NewLocalRecipeRepository(filePath string) (*LocalRecipeRepository, error) {
	store, err := newLocalStore[cooking.Recipe](filePath)
	if err != nil {
		return nil, err
	}
	return &LocalRecipeRepository{store: store}, nil
}

func (r *LocalRecipeRepository) List() ([]cooking.Recipe, error) {
	items, err := r.store.all()
	if err != nil {
		return nil, err
	}

	recipes := make([]cooking.Recipe, 0, len(items))
	for _, recipe := range items {
		recipes = append(recipes, recipe)
	}
	sort.SliceStable(recipes, func(i, j int) bool {
		if recipes[i].Name != recipes[j].Name {
			return recipes[i].Name < recipes[j].Name
		}
		return recipes[i].Id < recipes[j].Id
	})

	return recipes, nil
}

func (r *LocalRecipeRepository) Get(id string) (cooking.Recipe, error) {
	recipe, exists, err := r.store.get(id)
	if err != nil {
		return cooking.Recipe{}, err
	}
	if !exists {
		return cooking.Recipe{}, fmt.Errorf("%w: %s", cooking.ErrRecipeNotFound, id)
	}
	return recipe, nil
}

func (r *LocalRecipeRepository) CreateOrUpdate(recipe cooking.Recipe) error {
	return r.store.put(recipe.Id, recipe)
}

func (r *LocalRecipeRepository) Delete(id string) error {
	return r.store.remove(id)
}
