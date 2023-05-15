package feeds

import (
	"encoding/json"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/filecacher"
	"whitford.io/radiofooty/internal/interchange"
)

func getBBCMatches() []interchange.MergedMatch {
	var matches = []interchange.MergedMatch{}
	longFormat := "2006-01-02T15:04:05Z"
	baseUrl := "https://rms.api.bbc.co.uk/v2/experience/inline/schedules/bbc_radio_five_live/"
	urls := []string{}
	urlTime := "2006-01-02"
	loc, _ := time.LoadLocation("Europe/London")
	start := time.Now()
	for i := 0; i < 8; i++ {
		t := start.AddDate(0, 0, i)
		urls = append(urls, baseUrl+t.Format(urlTime))
	}

	var bbcFeed interchange.BBCFeed
	for _, url := range urls {
		body, err := filecacher.GetUrl(url)
		if err != nil {
			panic(err)
		}
		json.Unmarshal(body, &bbcFeed)

		for _, data := range bbcFeed.Data {
			for _, prog := range data.Data {
				if strings.HasPrefix(prog.Title.Secondary, "Women") {
					continue
				}
				if isLeagueGame(prog.Title) {
					start, _ := time.Parse(longFormat, prog.Start)
					start = start.In(loc)
					clock := start.Format(timeLayout)
					date := start.Format(niceDate)
					m := interchange.MergedMatch{Time: clock, Date: date, Station: "BBC Radio 5", Datetime: start.Format(time.RFC3339), Title: prog.Title.Tertiary, Competition: prog.Title.Secondary}
					matches = append(matches, m)
				}
			}
		}
	}

	return matches
}

func isLeagueGame(title interchange.BBCTitles) bool {
	return (title.Primary == "5 Live Sport") &&
		(strings.Contains(title.Secondary, "Football")) &&
		strings.Contains(title.Tertiary, " v ")
}
