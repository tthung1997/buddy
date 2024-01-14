package bgg

type User struct {
	Id               string `xml:"id,attr"`
	Username         string `xml:"name,attr"`
	FistName         string `xml:"firstname"`
	LastName         string `xml:"lastname"`
	Avatar           string `xml:"avatarlink"`
	YearRegistered   int    `xml:"yearregistered"`
	LastLogin        string `xml:"lastlogin"`
	StateOrProvince  string `xml:"stateorprovince"`
	Country          string `xml:"country"`
	WebAddress       string `xml:"webaddress"`
	XboxAccount      string `xml:"xboxaccount"`
	WiiAccount       string `xml:"wiiaccount"`
	PsnAccount       string `xml:"psnaccount"`
	BattleNetAccount string `xml:"battlenetaccount"`
	SteamAccount     string `xml:"steamaccount"`
	TraderRating     int    `xml:"traderating"`
}
