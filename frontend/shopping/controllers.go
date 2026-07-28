package shopping

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const staticDir = "frontend/shopping/static"
const dbDir = "frontend/shopping/.db"
const inventoryFile = dbDir + "/inventory.json"
const shoppingListFile = dbDir + "/shopping_list.json"

var indexTmpl *template.Template = template.Must(template.ParseFiles(staticDir + "/index.html"))

type InventoryItem struct {
	Name  string `json:"name"`
	Date  string `json:"date"`
	Store string `json:"store"`
}

type ShoppingItem struct {
	Name  string `json:"name"`
	Store string `json:"store"`
}

func Index(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, nil)
}

func InventoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := loadItems[InventoryItem](inventoryFile)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeItems(w, items)
	case http.MethodPost:
		var items []InventoryItem
		json.NewDecoder(r.Body).Decode(&items)
		if err := saveItems(inventoryFile, items); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func ShoppingListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		items, err := loadItems[ShoppingItem](shoppingListFile)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		writeItems(w, items)
	case http.MethodPost:
		var items []ShoppingItem
		json.NewDecoder(r.Body).Decode(&items)
		if err := saveItems(shoppingListFile, items); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

// AppendItems queues items on the shopping list and returns how many were
// added. Items already on the list are skipped so another module can push the
// same request twice without creating duplicates.
func AppendItems(items []ShoppingItem) (int, error) {
	list, err := loadItems[ShoppingItem](shoppingListFile)
	if err != nil {
		return 0, err
	}

	queued := map[string]bool{}
	for _, item := range list {
		queued[strings.ToLower(strings.TrimSpace(item.Name))] = true
	}

	added := 0
	for _, item := range items {
		item.Name = strings.TrimSpace(item.Name)
		if item.Name == "" || queued[strings.ToLower(item.Name)] {
			continue
		}
		list = append(list, item)
		queued[strings.ToLower(item.Name)] = true
		added++
	}

	if added == 0 {
		return 0, nil
	}
	return added, saveItems(shoppingListFile, list)
}

func loadItems[T any](filePath string) ([]T, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []T{}, nil
		}
		return nil, err
	}

	items := []T{}
	if len(data) == 0 {
		return items, nil
	}
	if err := json.Unmarshal(data, &items); err != nil {
		return nil, err
	}
	return items, nil
}

func saveItems[T any](filePath string, items []T) error {
	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return err
	}
	if items == nil {
		items = []T{}
	}

	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, data, os.ModePerm)
}

func writeItems[T any](w http.ResponseWriter, items []T) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}
