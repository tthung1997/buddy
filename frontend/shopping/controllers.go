package shopping

import (
	"html/template"
	"net/http"
)

const staticDir = "frontend/shopping/static"

var indexTmpl *template.Template = template.Must(template.ParseFiles(staticDir + "/index.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, nil)
}
