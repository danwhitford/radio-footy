package feeds

import (
	"fmt"
	"log"
	"regexp"
	"time"
)

const niceDate = "Monday, Jan 2"
const timeLayout = "15:04"
const CalTimeString = "20060102T150405Z"

type MatchDay struct {
	DateTime time.Time
	Matches  []Listing
}

type Match struct {
	Datetime    time.Time
	HomeTeam    string
	AwayTeam    string
	Competition string
}

type Broadcast struct {
	Match
	Station Station
}

type Listing struct {
	Match
	Stations []Station
}

func (m Match) Title() string {
	if m.Competition == "NFL" {
		return fmt.Sprintf("%s @ %s", m.AwayTeam, m.HomeTeam)
	}
	return fmt.Sprintf("%s v %s", m.HomeTeam, m.AwayTeam)
}

func (m Match) Time() string {
	return m.Datetime.Format(timeLayout)
}

func (match Match) RollUpHash() string {
	return fmt.Sprintf("%s%v%s%s",
		match.Competition,
		match.Datetime.Format(time.DateOnly),
		match.HomeTeam,
		match.AwayTeam,
	)
}

func (m MatchDay) NiceDate() string {
	return m.DateTime.Format(niceDate)
}

func (m MatchDay) DateOnly() string {
	return m.DateTime.Format(time.DateOnly)
}

func (l Listing) GameHash() string {
	s := fmt.Sprintf("%s%s", l.Datetime.Format(time.RFC3339), l.Title())
	r := regexp.MustCompile("[^0-9a-zA-Z]")
	s = r.ReplaceAllString(s, "")
	return s
}

type Station struct {
	Name string
	Rank int
}

func (stn Station) String() string {
	return stn.Name
}

var SkySports = Station{"Sky Sports", 0}
var TNTSports = Station{"TNT Sports", 10}
var BBCOne = Station{"BBC One", 20}
var BBCTwo = Station{"BBC Two", 30}
var ITV1 = Station{"ITV1", 40}
var ITV4 = Station{"ITV4", 44}
var ChannelFour = Station{"Channel 4", 50}
var Talksport = Station{"talkSPORT", 60}
var Talksport2 = Station{"talkSPORT2", 70}
var Radio5 = Station{"Radio 5 Live", 80}
var Radio5Extra = Station{"Radio 5 Sports Extra", 90}
var BlankStation = Station{"", 9999}

func StationFromString(name string) Station {
	for _, station := range []Station{
		SkySports,
		TNTSports,
		BBCOne,
		BBCTwo,
		ITV1,
		ITV4,
		ChannelFour,
		Talksport,
		Talksport2,
		Radio5,
		Radio5Extra,
	} {
		if name == station.Name {
			return station
		}
	}
	log.Printf("station not found: '%s'\n", name)
	return Station{name, 9999}
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
	Location []Station
}
