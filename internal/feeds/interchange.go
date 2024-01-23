package feeds

import (
	"fmt"
	"time"
)

type MatchDay struct {
	NiceDate string
	DateTime time.Time
	Matches  []Listing
}

type Match struct {
	Time        string
	Date        string
	Datetime    string
	HomeTeam    string
	AwayTeam    string
	Competition string
}

type Broadcast struct {
	Match
	Station string
}

type Listing struct {
	Match
	Stations []string
}

func (m Match) Title() string {
	if m.Competition == "NFL" {
		return fmt.Sprintf("%s @ %s", m.AwayTeam, m.HomeTeam)
	}
	return fmt.Sprintf("%s v %s", m.HomeTeam, m.AwayTeam)
}

type MatchRadioEvent struct {
	Station string
	Time    string
	Date    string
}

type TSGame struct {
	Livefeed []TSLiveFeed `json:"livefeed"`
	Sport    string
	Date     string
	HomeTeam string
	AwayTeam string
	League   string
	Title    string
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
	Title    BBCTitles   `json:"titles"`
	Start    string      `json:"start"`
	Network  BBCNetwork  `json:"network"`
	Synopses BBCSynopses `json:"synopses"`
}

type BBCSynopses struct {
	Short string `json:"short"`
}

type BBCNetwork struct {
	ShortTitle string `json:"short_title"`
}

type BBCTitles struct {
	Primary   string `json:"primary"`
	Secondary string `json:"secondary"`
	Tertiary  string `json:"tertiary"`
}

type CalEvent struct {
	Uid      string
	DtStart  string
	Summary  string
	Location []string
}

const CalTimeString = "20060102T150405Z"
