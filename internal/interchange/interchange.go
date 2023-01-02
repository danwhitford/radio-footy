package interchange

type Merged = []MergedMatchDay

type MergedMatchDay struct {
	Date    string        `json:"date"`
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
	Sport    string
	Date     string
	HomeTeam string
	AwayTeam string
	League   string
}

type TSLiveFeed struct {
	Feedname string `json:"feedname"`
}

type BBCFeed struct {
	Data []BBCFeedData `json:"data"`
}

type BBCFeedData struct {
	Data []BBCProgramData `json:"data"`
}

type BBCProgramData struct {
	Title BBCTitles `json:"titles"`
	Start string    `json:"start"`
}

type BBCTitles struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Tertiary  string `json:"tertiary"`
}

type CalEvent struct {
	Uid         string
	DtStart     string
	Summary     string
	Location    string
}

const CalTimeString = "20060102T150405Z"
