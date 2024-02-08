package feeds

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/urlgetter"
)

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

func getBBCMatches(getter urlgetter.UrlGetter) ([]Broadcast, error) {
	var matches = make([]Broadcast, 0)
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
		err = json.Unmarshal(body, &bbcFeed)
		if err != nil {
			return nil, err
		}

		merged := bbcDayToMatches(bbcFeed)
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

func isSixNations(title BBCTitles) bool {
	return title.Primary == "Six Nations 2024" &&
		strings.Contains(title.Secondary, " v ")
}

func bbcDayToMatches(bbcFeed BBCFeed) []Broadcast {
	matches := make([]Broadcast, 0)

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
				teams := strings.Split(prog.Title.Tertiary, " v ")
				m := NewSantisedMatch(
					start,
					teams[0],
					teams[1],
					prog.Title.Secondary,
				)

				matches = append(matches, Broadcast{m, StationFromString(prog.Network.ShortTitle)})
			} else if isSixNations(prog.Title) {
				start, err := time.Parse(longFormat, prog.Start)
				if err != nil {
					panic(err)
				}
				start = start.In(loc)
				teams := strings.Split(prog.Title.Secondary, " v ")
				m := NewSantisedMatch(
					start,
					teams[0],
					teams[1],
					prog.Title.Primary,
				)

				matches = append(matches, Broadcast{m, StationFromString(prog.Network.ShortTitle)})
			}
		}
	}

	return matches
}
