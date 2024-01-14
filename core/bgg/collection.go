package bgg

type CollectionFilter struct {
	Username       string
	Subtype        string
	ExcludeSubtype string
	Own            bool
	Rated          bool
	Played         bool
	Trade          bool
	Want           bool
	Wishlist       bool
	Preordered     bool
	PrevOwned      bool
}

type CollectionItem struct {
	Id            string `xml:"objectid,attr"`
	Subtype       string `xml:"subtype,attr"`
	CollectionId  string `xml:"collid,attr"`
	Name          string `xml:"name"`
	YearPublished string `xml:"yearpublished"`
	Image         string `xml:"image"`
	Thumbnail     string `xml:"thumbnail"`
	NumPlays      string `xml:"numplays"`
	Status        struct {
		Own          string `xml:"own,attr"`
		PrevOwned    string `xml:"prevowned,attr"`
		ForTrade     string `xml:"fortrade,attr"`
		Want         string `xml:"want,attr"`
		WantToPlay   string `xml:"wanttoplay,attr"`
		WantToBuy    string `xml:"wanttobuy,attr"`
		Wishlist     string `xml:"wishlist,attr"`
		Preordered   string `xml:"preordered,attr"`
		LastModified string `xml:"lastmodified,attr"`
	} `xml:"status"`
	Comments      string `xml:"comment"`
	ConditionText string `xml:"conditiontext"`
}

type Collection struct {
	Items []CollectionItem `xml:"item"`
}
