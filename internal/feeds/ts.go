package feeds

import (
	"encoding/json"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/urlgetter"
)

func getTalkSportMatches(getter urlgetter.UrlGetter) ([]Match, error) {
	url := "https://talksport.com/wp-json/talksport/v2/talksport-live/commentary"

	body, err := getter.GetUrl(url)
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

func tsFeedToMatches(tsFeed []TSGame) []Match {
	matches := make([]Match, 0)
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
		if m.League == "" {
			continue
		}
		if strings.HasPrefix(m.Title, "Women") ||
			strings.Contains(m.Title, "Women") {
			continue
		}

		t, _ := time.ParseInLocation(longForm, m.Date, loc)
		displayDate := t.Format(niceDate)
		displayTime := t.Format(timeLayout)
		datetime := t.Format(time.RFC3339)
		m := Match{
			Time: displayTime, 
			Date: displayDate, 
			Stations: []string{feedname}, 
			Datetime: datetime, 
			HomeTeam: m.HomeTeam,
			AwayTeam: m.AwayTeam,
			Competition: m.League}
		matches = append(matches, m)
	}

	return matches
}
