package boardgames

import "github.com/tthung1997/buddy/core/bgg"

type IndexPageData struct {
	Error      error
	Filter     bgg.CollectionFilter
	Collection bgg.Collection
}

type PickPageData struct {
	Error error
	Items []bgg.CollectionItem
}
