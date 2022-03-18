package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"whitford.io/radiofooty/internal/interchange"
)

func main() {
	url := "https://talksport.com/wp-json/talksport/v2/talksport-live/commentary"
	resp, _ := http.Get(url)
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tsFeed interchange.TSFeed
	json.Unmarshal(body, &tsFeed)

	var mergedFeed interchange.Merged
	var matches []interchange.MergedMatch
	longForm := "2006-01-02 15:04:05"
	niceDate := "Monday, Jan 2nd"
	timeLayout := "15:04"
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
		t, _ := time.ParseInLocation(longForm, m.Date, loc)
		title := fmt.Sprintf("%s v %s", m.HomeTeam, m.AwayTeam)
		displayDate := t.Format(niceDate)
		displayTime := t.Format(timeLayout)
		datetime := t.Format(time.RFC3339)
		m := interchange.MergedMatch{Time: displayTime, Date: displayDate, Station: feedname, Datetime: datetime, Title: title, Competition: m.League}
		matches = append(matches, m)		
	}

	matchesRollup := make(map[string][]interchange.MergedMatch)
	for _, match := range matches {
		val, prs := matchesRollup[match.Date]
		if prs {
			val = append(val, match)
			matchesRollup[match.Date] = val
		} else {
			matchesRollup[match.Date] = []interchange.MergedMatch{match}
		}
	}
	
	for k, v :=  range matchesRollup {
		md := interchange.MergedMatchDay{Date: k, Matches: v}
		mergedFeed = append(mergedFeed, md)
	}

	output, _ := json.Marshal(mergedFeed)
	fmt.Println(string(output))
}
