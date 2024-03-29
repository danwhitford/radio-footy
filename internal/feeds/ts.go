package feeds

import (
	"encoding/json"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/urlgetter"
)

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

type talkSportGetter struct {
	urlgetter urlgetter.UrlGetter
}

func (tsg talkSportGetter) getMatches() ([]Broadcast, error) {
	url := "https://talksport.com/wp-json/talksport/v2/talksport-live/commentary"

	body, err := tsg.urlgetter.GetUrl(url)
	if err != nil {
		return nil, err
	}
	var tsFeed []TSGame
	err = json.Unmarshal(body, &tsFeed)
	if err != nil {
		return nil, err
	}

	return tsFeedToMatches(tsFeed), nil
}

func tsFeedToMatches(tsFeed []TSGame) []Broadcast {
	matches := make([]Broadcast, 0)
	longForm := "2006-01-02 15:04:05"
	loc, _ := time.LoadLocation("Europe/London")

	for _, m := range tsFeed {
		var feedname string
		for _, feed := range m.Livefeed {
			if feed.Feedname == "talkSPORT" {
				feedname = "talkSPORT"
			} else if feed.Feedname == "talkSPORT2" {
				feedname = "talkSPORT2"
			}
		}
		if feedname == "" {
			continue
		}
		if m.League == "" || m.League == "International Friendlies" {
			continue
		}

		if strings.HasPrefix(m.Title, "Women") ||
			strings.Contains(m.Title, "Women") {
			continue
		}

		t, _ := time.ParseInLocation(longForm, m.Date, loc)
		datetime := t
		m := NewSantisedMatch(
			datetime,
			m.HomeTeam,
			m.AwayTeam,
			m.League,
		)
		matches = append(matches, Broadcast{m, StationFromString(feedname)})
	}

	return matches
}
