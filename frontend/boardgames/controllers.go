package boardgames

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
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
