package boardgames

import (
	"html/template"
	"net/http"

	"github.com/tthung1997/buddy/core/bgg"
)

const boardgamesStaticDir = "frontend/boardgames/static"

var indexTmpl *template.Template = template.Must(template.ParseFiles(boardgamesStaticDir + "/index.html"))
var bggClient = bgg.NewClient(*bgg.DefaultClientConfig())

func Index(w http.ResponseWriter, r *http.Request) {
	collection, err := bggClient.GetCollection(
		bgg.CollectionFilter{
			Username: "tthung1997",
		},
	)
	if err != nil {
		indexTmpl.Execute(w, IndexPageData{
			Error: err,
		})
		return
	}

	indexTmpl.Execute(w, IndexPageData{
		Collection: *collection,
	})
}
