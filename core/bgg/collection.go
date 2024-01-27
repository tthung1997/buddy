package bgg

type CollectionFilter struct {
	Username       string
	Subtype        string
	ExcludeSubtype string
	Own            string
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
	YearPublished int    `xml:"yearpublished"`
	Image         string `xml:"image"`
	Thumbnail     string `xml:"thumbnail"`
	NumPlays      int    `xml:"numplays"`
	Status        struct {
		Own          int    `xml:"own,attr"`
		PrevOwned    int    `xml:"prevowned,attr"`
		ForTrade     int    `xml:"fortrade,attr"`
		Want         int    `xml:"want,attr"`
		WantToPlay   int    `xml:"wanttoplay,attr"`
		WantToBuy    int    `xml:"wanttobuy,attr"`
		Wishlist     int    `xml:"wishlist,attr"`
		Preordered   int    `xml:"preordered,attr"`
		LastModified string `xml:"lastmodified,attr"`
	} `xml:"status"`
	Comments      string `xml:"comment"`
	ConditionText string `xml:"conditiontext"`
}

type Collection struct {
	Items []CollectionItem `xml:"item"`
}
