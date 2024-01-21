package feeds

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/filecacher"
)

func getBBCMatches(getter filecacher.Getter) ([]MergedMatch, error) {
	var matches = []MergedMatch{}
	baseUrls := []string{
		"https://rms.api.bbc.co.uk/v2/experience/inline/schedules/bbc_radio_five_live/",
		"https://rms.api.bbc.co.uk/v2/experience/inline/schedules/bbc_radio_five_live_sports_extra/",
	}
	urls := []string{}
	urlTime := "2006-01-02"
	start := time.Now()
	for i := 0; i < 8; i++ {
		t := start.AddDate(0, 0, i)
		for _, baseUrl := range baseUrls {
			urls = append(urls, baseUrl+t.Format(urlTime))
		}
	}

	var bbcFeed BBCFeed
	for _, url := range urls {
		body, err := getter.GetUrl(url)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(body, &bbcFeed)

		merged := bbcDayToMergedMatches(bbcFeed)
		matches = append(matches, merged...)
	}

	return matches, nil
}

func isAMatch(title BBCTitles) bool {
	return (title.Primary == "5 Live Sport") &&
		(strings.Contains(title.Secondary, "Football")) &&
		strings.Contains(title.Tertiary, " v ") &&
		!strings.HasPrefix(title.Tertiary, "Prematch")
}

func bbcDayToMergedMatches(bbcFeed BBCFeed) []MergedMatch {
	matches := make([]MergedMatch, 0)

	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		log.Fatalf("error loading location: %v", err)
	}
	longFormat := "2006-01-02T15:04:05Z"

	for _, data := range bbcFeed.Data {
		for _, prog := range data.Data {
			if strings.HasPrefix(prog.Title.Secondary, "Women") ||
				strings.Contains(prog.Title.Tertiary, "Women") ||
				strings.Contains(prog.Synopses.Short, "Women") {
				continue
			}
			if isAMatch(prog.Title) {
				start, err := time.Parse(longFormat, prog.Start)
				if err != nil {
					panic(err)
				}
				start = start.In(loc)
				clock := start.Format(timeLayout)
				date := start.Format(niceDate)
				m := MergedMatch{
					Time:        clock,
					Date:        date,
					Stations:    []string{"BBC Radio 5"},
					Datetime:    start.Format(time.RFC3339),
					Title:       prog.Title.Tertiary,
					Competition: prog.Title.Secondary,
				}

				matches = append(matches, m)
			}
		}
	}

	return matches
}
