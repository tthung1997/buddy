package boardgames

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/patrickmn/go-cache"

	"github.com/tthung1997/buddy/core/bgg"
)

const boardgamesStaticDir = "frontend/boardgames/static"

var indexTmpl *template.Template = template.Must(template.ParseFiles(boardgamesStaticDir + "/index.html"))
var bggClient = bgg.NewClient(*bgg.DefaultClientConfig())
var cachedCollections = cache.New(1*time.Minute, 1*time.Minute)

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

	return filter
}

func Index(w http.ResponseWriter, r *http.Request) {
	// Get query params
	params := r.URL.Query()
	filter := getFilter(params)

	// Get collection
	collectionFilter := bgg.CollectionFilter{
		Username:       filter.Username,
		ExcludeSubtype: filter.ExcludeSubtype,
	}
	jsonBytes, err := json.Marshal(collectionFilter)
	if err != nil {
		indexTmpl.Execute(w, IndexPageData{
			Error: err,
		})
		return
	}
	filterString := string(jsonBytes)
	var collection *bgg.Collection
	if x, found := cachedCollections.Get(filterString); found {
		collection = x.(*bgg.Collection)
	} else {
		var err error

		collection, err = bggClient.GetCollection(collectionFilter)
		if err != nil {
			indexTmpl.Execute(w, IndexPageData{
				Error: err,
			})
			return
		}
		cachedCollections.Set(filterString, collection, cache.DefaultExpiration)
	}

	indexTmpl.Execute(w, IndexPageData{
		Filter:     filter,
		Collection: *collection,
	})
}
