package repository

import (
	"fmt"
	"sort"

	"github.com/tthung1997/buddy/core/cooking"
)

type LocalPantryRepository struct {
	store *localStore[cooking.PantryItem]
}

func NewLocalPantryRepository(filePath string) (*LocalPantryRepository, error) {
	store, err := newLocalStore[cooking.PantryItem](filePath)
	if err != nil {
		return nil, err
	}
	return &LocalPantryRepository{store: store}, nil
}

func (r *LocalPantryRepository) List() ([]cooking.PantryItem, error) {
	items, err := r.store.all()
	if err != nil {
		return nil, err
	}

	pantry := make([]cooking.PantryItem, 0, len(items))
	for _, item := range items {
		pantry = append(pantry, item)
	}
	sort.SliceStable(pantry, func(i, j int) bool {
		return pantry[i].IngredientId < pantry[j].IngredientId
	})

	return pantry, nil
}

func (r *LocalPantryRepository) Get(ingredientId string) (cooking.PantryItem, error) {
	item, exists, err := r.store.get(ingredientId)
	if err != nil {
		return cooking.PantryItem{}, err
	}
	if !exists {
		return cooking.PantryItem{}, fmt.Errorf("%w: %s", cooking.ErrPantryItemNotFound, ingredientId)
	}
	return item, nil
}

func (r *LocalPantryRepository) CreateOrUpdate(item cooking.PantryItem) error {
	return r.store.put(item.IngredientId, item)
}

func (r *LocalPantryRepository) Delete(ingredientId string) error {
	return r.store.remove(ingredientId)
}
