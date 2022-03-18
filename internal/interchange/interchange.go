package interchange

type Merged = []MergedMatchDay

type MergedMatchDay struct {
	Date    string `json:"date"`
	Matches []MergedMatch `json:"matches"`
}

type MergedMatch struct {
	Time        string `json:"time"`
	Date        string `json:"date"`
	Station     string `json:"station"`
	Datetime    string `json:"datetime"`
	Title       string `json:"title"`
	Competition string `json:"competition"`
}

type TSFeed = []TSGames

type TSGames struct {
	Livefeed []TSLiveFeed `json:"livefeed"`
	Sport string
	Date string
	HomeTeam string
	AwayTeam string
	League string
}

type TSLiveFeed struct {
	Feedname string `json:"feedname"`
}