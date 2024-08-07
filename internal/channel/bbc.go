package channel

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/broadcast"
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

type BbcMatchGetter struct {
	Urlgetter urlgetter.UrlGetter
}

func (bbc BbcMatchGetter) GetMatches() ([]broadcast.Broadcast, error) {
	matches := make([]broadcast.Broadcast, 0)
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
		body, err := bbc.Urlgetter.GetUrl(url)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(body, &bbcFeed)
		if err != nil {
			return nil, err
		}

		merged, err := bbcDayToMatches(bbcFeed)
		if err != nil {
			return merged, err
		}
		matches = append(matches, merged...)
	}

	// matches = dedupeBbcMatches(matches)

	return matches, nil
}

func dedupeBbcMatches(matches []broadcast.Broadcast) []broadcast.Broadcast {
	rollUp := make(map[string][]broadcast.Broadcast)

	for _, b := range matches {
		rollUp[b.RollUpHash()] = append(rollUp[b.RollUpHash()], b)
	}

	unique := make([]broadcast.Broadcast, 0)
	for _, bb := range rollUp {
		sort.Slice(bb, func(i, j int) bool {
			return bb[i].Station.Rank < bb[j].Station.Rank
		})
		unique = append(unique, bb[0])
	}

	return unique
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

func isEuros(title BBCTitles) bool {
	return title.Primary == "Euro 2024" && strings.Contains(title.Secondary, " v ") && !strings.Contains(title.Secondary, "Prematch")
}

func isCricket(title BBCTitles) bool {
	return title.Primary == "Test Match Special" && strings.Contains(title.Secondary, " v ")
}

func bbcDayToMatches(bbcFeed BBCFeed) ([]broadcast.Broadcast, error) {
	matches := make([]broadcast.Broadcast, 0)

	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		return matches, fmt.Errorf("error loading location: %v", err)
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
				m := broadcast.NewSantisedMatch(
					start,
					teams[0],
					teams[1],
					prog.Title.Secondary,
				)

				matches = append(matches, broadcast.Broadcast{Match: m, Station: broadcast.StationFromString(prog.Network.ShortTitle)})
			} else if isSixNations(prog.Title) {
				start, err := time.Parse(longFormat, prog.Start)
				if err != nil {
					panic(err)
				}
				start = start.In(loc)
				teams := strings.Split(prog.Title.Secondary, " v ")
				m := broadcast.NewSantisedMatch(
					start,
					teams[0],
					teams[1],
					prog.Title.Primary,
				)

				matches = append(matches, broadcast.Broadcast{Match: m, Station: broadcast.StationFromString(prog.Network.ShortTitle)})
			} else if isEuros(prog.Title) {
				start, err := time.Parse(longFormat, prog.Start)
				if err != nil {
					panic(err)
				}
				start = start.In(loc)
				teams := strings.Split(prog.Title.Secondary, " v ")
				m := broadcast.NewSantisedMatch(
					start,
					teams[0],
					teams[1],
					"UEFA EURO 2024",
				)

				matches = append(matches, broadcast.Broadcast{Match: m, Station: broadcast.Talksport})
				matches = append(matches, broadcast.Broadcast{Match: m, Station: broadcast.StationFromString(prog.Network.ShortTitle)})
			} else if isCricket(prog.Title) {
				split := strings.Split(prog.Title.Secondary, " - ")
				start, err := time.Parse(longFormat, prog.Start)
				if err != nil {
					panic(err)
				}
				start = start.In(loc)
				teams := strings.Split(split[0], " v ")
				m := broadcast.NewSantisedMatch(
					start,
					teams[0],
					teams[1],
					split[1],
				)

				matches = append(matches, broadcast.Broadcast{Match: m, Station: broadcast.StationFromString(prog.Network.ShortTitle)})
			}
		}
	}

	return matches, nil
}
