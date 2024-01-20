package feeds

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/filecacher"
)

func getTalkSportMatches() ([]MergedMatch, error) {
	url := "https://talksport.com/wp-json/talksport/v2/talksport-live/commentary"
	cacher := filecacher.CachedGetter{
		Getter: filecacher.NewHttpGetter(),
	}
	body, err := cacher.GetUrl(url)
	if err != nil {
		return nil, err
	}
	var tsFeed []TSGame
	err = json.Unmarshal(body, &tsFeed)
	if err != nil {
		return nil, err
	}

	return tsFeedToMergedMatches(tsFeed), nil
}

func tsFeedToMergedMatches(tsFeed []TSGame) []MergedMatch {
	matches := make([]MergedMatch, 0)
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
		title := fmt.Sprintf("%s v %s", m.HomeTeam, m.AwayTeam)
		displayDate := t.Format(niceDate)
		displayTime := t.Format(timeLayout)
		datetime := t.Format(time.RFC3339)
		m := MergedMatch{Time: displayTime, Date: displayDate, Stations: []string{feedname}, Datetime: datetime, Title: title, Competition: m.League}
		matches = append(matches, m)
	}

	return matches
}
