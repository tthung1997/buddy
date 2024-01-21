package home

import (
	"html/template"
	"net/http"
)

const staticDir = "frontend/home/static"

var indexTmpl *template.Template = template.Must(template.ParseFiles(staticDir + "/index.html"))

func Index(w http.ResponseWriter, r *http.Request) {
	indexTmpl.Execute(w, nil)
}
