package shopping

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"os"
)

const staticDir = "frontend/shopping/static"
const inventoryFile = "frontend/shopping/.db/inventory.json"
const shoppingListFile = "frontend/shopping/.db/shopping_list.json"

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
		f, err := os.Open(inventoryFile)
		if err == nil {
			defer f.Close()
			w.Header().Set("Content-Type", "application/json")
			io.Copy(w, f)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		var items []InventoryItem
		json.NewDecoder(r.Body).Decode(&items)
		b, _ := json.MarshalIndent(items, "", "  ")
		err := os.WriteFile(inventoryFile, b, os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func ShoppingListHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		f, err := os.Open(shoppingListFile)
		if err == nil {
			defer f.Close()
			w.Header().Set("Content-Type", "application/json")
			io.Copy(w, f)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	case http.MethodPost:
		var items []ShoppingItem
		json.NewDecoder(r.Body).Decode(&items)
		b, _ := json.MarshalIndent(items, "", "  ")
		err := os.WriteFile(shoppingListFile, b, os.ModePerm)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
