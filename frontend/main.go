package main

import (
	"log"
	"net/http"

	"github.com/tthung1997/buddy/frontend/boardgames"
	"github.com/tthung1997/buddy/frontend/home"
	"github.com/tthung1997/buddy/frontend/shopping"
)

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[Incoming] %v %v?%v", r.Method, r.URL.Path, r.URL.RawQuery)
		f(w, r)
	}
}

func main() {
	// no favicon
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	// home
	http.HandleFunc("/", logging(home.Index))

	// board games
	http.Handle("/boardgames/static/", http.StripPrefix("/boardgames/static/", http.FileServer(http.Dir("frontend/boardgames/static"))))
	http.HandleFunc("/boardgames", logging(boardgames.Index))

	// shopping
	http.HandleFunc("/shopping", logging(shopping.Index))

	log.Println("Listening on port 2210")
	err := http.ListenAndServe(":2210", nil)
	if err != nil {
		log.Fatal(err)
	}
}
