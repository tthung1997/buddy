package boardgames

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/tthung1997/buddy/core/bgg"
	coreRandom "github.com/tthung1997/buddy/core/random"
)

type BoardGamesController struct {
	BggClient         *bgg.Client
	Randomizer        coreRandom.IRandomizer
	CachedCollections *cache.Cache
}

func NewBoardGamesController(bggClient *bgg.Client, randomizer coreRandom.IRandomizer) *BoardGamesController {
	return &BoardGamesController{
		BggClient:         bggClient,
		Randomizer:        randomizer,
		CachedCollections: cache.New(1*time.Minute, 1*time.Minute),
	}
}

const boardgamesStaticDir = "frontend/boardgames/static"

var indexTmpl *template.Template = template.Must(template.ParseFiles(boardgamesStaticDir + "/index.html"))
var pickTmpl *template.Template = template.Must(template.ParseFiles(boardgamesStaticDir + "/pick.html"))

// --- Backed/Preordered Tracker ---
type BackedGame struct {
	Name      string
	Delivery  string  // month/year
	PricePaid float64 // Price paid for the game
	From      string  // Where the game was backed from
	Platform  string  // Platform for the game (e.g., Steam, Epic)
	Received  bool    // Whether the game has been received
}

type BackedPageData struct {
	Error   error
	Success string
	Items   []BackedGame
	AddName string
}

const backedDataFile = "frontend/boardgames/.db/backed_boardgames.json"

var backedTmpl *template.Template = template.Must(template.ParseFiles(boardgamesStaticDir + "/backed.html"))

// loadBackedGames loads the backed games from the JSON file
func loadBackedGames() ([]BackedGame, error) {
	f, err := os.Open(backedDataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []BackedGame{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var items []BackedGame
	err = json.NewDecoder(f).Decode(&items)
	if err != nil {
		return nil, err
	}
	return items, nil
}

// saveBackedGames saves the backed games to the JSON file
func saveBackedGames(items []BackedGame) error {
	// Ensure .db directory exists
	dbDir := "frontend/boardgames/.db"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return err
		}
	}
	f, err := os.Create(backedDataFile)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(items)
}

func getFilter(params url.Values) bgg.CollectionFilter {
	// Construct filter
	filter := bgg.CollectionFilter{}

	// Get username
	username := params.Get("username")
	if username == "" {
		username = "tthung1997"
	}
	filter.Username = username

	// Get excludeExpansion
	excludeExpansion := params.Get("excludeExpansion")
	if excludeExpansion == "on" {
		filter.ExcludeSubtype = "boardgameexpansion"
	}

	// Get own
	own := params.Get("own")
	if own == "" {
		own = "yes"
	}
	switch own {
	case "yes":
		filter.Own = "1"
	case "no":
		filter.Own = "0"
	}

	return filter
}

func (c *BoardGamesController) getCollection(filter bgg.CollectionFilter) (*bgg.Collection, error) {
	jsonBytes, err := json.Marshal(filter)
	if err != nil {
		return nil, err
	}
	filterString := string(jsonBytes)
	var collection *bgg.Collection
	if x, found := c.CachedCollections.Get(filterString); found {
		collection = x.(*bgg.Collection)
	} else {
		var err error

		collection, err = c.BggClient.GetCollection(filter)
		if err != nil {
			return nil, err
		}

		c.CachedCollections.Set(filterString, collection, cache.DefaultExpiration)
	}

	return collection, nil
}

func (c *BoardGamesController) Index(w http.ResponseWriter, r *http.Request) {
	// Get query params
	params := r.URL.Query()
	filter := getFilter(params)

	// Get collection
	collectionFilter := bgg.CollectionFilter{
		Username:       filter.Username,
		ExcludeSubtype: filter.ExcludeSubtype,
	}
	collection, err := c.getCollection(collectionFilter)
	if err != nil {
		indexTmpl.Execute(w, IndexPageData{
			Error: err,
		})
		return
	}

	indexTmpl.Execute(w, IndexPageData{
		Filter:     filter,
		Collection: *collection,
	})
}

func (c *BoardGamesController) Pick(w http.ResponseWriter, r *http.Request) {
	// Get query
	params := r.URL.Query()
	filter := getFilter(params)
	countParam := params.Get("count")
	if countParam == "" {
		countParam = "1"
	}
	count, err := strconv.Atoi(countParam)
	if err != nil {
		pickTmpl.Execute(w, PickPageData{
			Error: err,
		})
		return
	}
	prioritizeLessPlayed := params.Get("prioritizeLessPlayed") == "on"

	// Get collection
	collection, err := c.getCollection(filter)
	if err != nil {
		pickTmpl.Execute(w, PickPageData{
			Error: err,
		})
		return
	}

	// Get max play count
	maxPlayCount := 0
	for _, item := range collection.Items {
		if item.NumPlays > maxPlayCount {
			maxPlayCount = item.NumPlays
		}
	}

	// Convert collection to choices with equal weight
	choices := make([]coreRandom.Choice, len(collection.Items))
	for i, item := range collection.Items {
		var weight int32 = 1
		if prioritizeLessPlayed {
			weight = int32(maxPlayCount - item.NumPlays + 1)
		}

		choices[i] = coreRandom.Choice{
			Value:  item.Id,
			Weight: weight,
		}
	}

	// Pick random choices
	pickedChoices := c.Randomizer.GetChoice(choices, count)

	// Extract picked choices' names and return as JSON
	pickedItems := make([]bgg.CollectionItem, len(pickedChoices))
	for i, choice := range pickedChoices {
		for _, item := range collection.Items {
			if choice.Value == item.Id {
				pickedItems[i] = item
				break
			}
		}
	}

	pickTmpl.Execute(w, PickPageData{
		Items: pickedItems,
	})
}

// Backed handles the backed games page
func (c *BoardGamesController) Backed(w http.ResponseWriter, r *http.Request) {
	items, _ := loadBackedGames()
	data := BackedPageData{Items: items}
	if name := r.URL.Query().Get("name"); name != "" {
		data.AddName = name
	}
	backedTmpl.Execute(w, data)
}

// Add new backed game
func (c *BoardGamesController) BackedAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/boardgames/backed", http.StatusSeeOther)
		return
	}
	items, _ := loadBackedGames()
	price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
	game := BackedGame{
		Name:      r.FormValue("name"),
		Delivery:  r.FormValue("delivery"),
		PricePaid: price,
		From:      r.FormValue("from"),
		Platform:  r.FormValue("platform"),
		Received:  false,
	}
	items = append(items, game)
	if err := saveBackedGames(items); err != nil {
		http.Error(w, "Failed to save backed games: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/boardgames/backed", http.StatusSeeOther)
}

// Mark as received
func (c *BoardGamesController) BackedReceive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/boardgames/backed", http.StatusSeeOther)
		return
	}
	items, _ := loadBackedGames()
	idx, _ := strconv.Atoi(r.FormValue("idx"))
	if idx >= 0 && idx < len(items) {
		items[idx].Received = true
		saveBackedGames(items)
	}
	http.Redirect(w, r, "/boardgames/backed", http.StatusSeeOther)
}

// Edit backed game (GET: show form, POST: save changes)
func (c *BoardGamesController) BackedEdit(w http.ResponseWriter, r *http.Request) {
	items, _ := loadBackedGames()
	idx, _ := strconv.Atoi(r.FormValue("idx"))
	if r.Method == http.MethodGet {
		if idx < 0 || idx >= len(items) {
			http.Redirect(w, r, "/boardgames/backed", http.StatusSeeOther)
			return
		}
		game := items[idx]
		// Render edit form (reuse backed.html for simplicity, pre-fill .AddName and other fields)
		data := BackedPageData{Items: items, AddName: game.Name}
		backedTmpl.Execute(w, data)
		return
	}
	// POST: update
	if idx >= 0 && idx < len(items) {
		price, _ := strconv.ParseFloat(r.FormValue("price"), 64)
		items[idx] = BackedGame{
			Name:      r.FormValue("name"),
			Delivery:  r.FormValue("delivery"),
			PricePaid: price,
			From:      r.FormValue("from"),
			Platform:  r.FormValue("platform"),
			Received:  items[idx].Received,
		}
		saveBackedGames(items)
	}
	http.Redirect(w, r, "/boardgames/backed", http.StatusSeeOther)
}
