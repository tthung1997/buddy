package main

import (
	"net/http"

	"github.com/tthung1997/buddy/frontend/boardgames"
	"github.com/tthung1997/buddy/frontend/home"
	"github.com/tthung1997/buddy/frontend/shopping"
)

func main() {
	http.HandleFunc("/", home.Index)

	// board games
	http.Handle("/boardgames/static/", http.StripPrefix("/boardgames/static/", http.FileServer(http.Dir("frontend/boardgames/static"))))
	http.HandleFunc("/boardgames", boardgames.Index)

	// shopping
	http.HandleFunc("/shopping", shopping.Index)

	http.ListenAndServe(":2210", nil)
}
