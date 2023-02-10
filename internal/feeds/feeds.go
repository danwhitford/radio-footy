package feeds

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/filecacher"
	"whitford.io/radiofooty/internal/interchange"
)

const niceDate = "Monday, Jan 2"
const timeLayout = "15:04"

func getTalkSportMatches() []interchange.MergedMatch {
	url := "https://talksport.com/wp-json/talksport/v2/talksport-live/commentary"
	body, _ := filecacher.GetUrl(url)
	var tsFeed interchange.TSFeed
	json.Unmarshal(body, &tsFeed)

	var matches []interchange.MergedMatch
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
		t, _ := time.ParseInLocation(longForm, m.Date, loc)
		title := fmt.Sprintf("%s v %s", m.HomeTeam, m.AwayTeam)
		displayDate := t.Format(niceDate)
		displayTime := t.Format(timeLayout)
		datetime := t.Format(time.RFC3339)
		m := interchange.MergedMatch{Time: displayTime, Date: displayDate, Station: feedname, Datetime: datetime, Title: title, Competition: m.League}
		matches = append(matches, m)
	}

	return matches
}

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
		body, _ := filecacher.GetUrl(url)
		json.Unmarshal(body, &bbcFeed)

		for _, data := range bbcFeed.Data {
			for _, prog := range data.Data {
				if isWorldCup(prog.Title) {
					start, _ := time.Parse(longFormat, prog.Start)
					start = start.In(loc)
					clock := start.Format(timeLayout)
					date := start.Format(niceDate)
					m := interchange.MergedMatch{Time: clock, Date: date, Station: "BBC Radio 5", Datetime: start.Format(time.RFC3339), Title: prog.Title.Tertiary, Competition: "World Cup"}
					matches = append(matches, m)
				} else if isLeagueGame(prog.Title) {
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

func isWorldCup(title interchange.BBCTitles) bool {
	return title.Primary == "World Cup" && strings.Contains(title.Tertiary, " v ")
}

func isLeagueGame(title interchange.BBCTitles) bool {
	return (title.Primary == "5 Live Sport") &&
		(strings.Contains(title.Secondary, "Football")) &&
		strings.Contains(title.Tertiary, " v ")
}

func GetMergedMatches() []interchange.MergedMatchDay {
	var mergedFeed interchange.Merged
	var matches []interchange.MergedMatch

	matches = append(matches, getTalkSportMatches()...)
	matches = append(matches, getBBCMatches()...)

	m := make([]interchange.MergedMatch, 0)
	for _, match := range matches {
		if shouldSkip(match.Competition) {
			continue
		}
		m = append(m, match)
	}
	matches = m

	for i := range matches {
		teams := strings.Split(matches[i].Title, " v ")
		newTitle := fmt.Sprintf("%s v %s", mapTeamName(teams[0]), mapTeamName(teams[1]))
		newComp := mapCompName(matches[i].Competition)
		matches[i].Title = newTitle
		matches[i].Competition = newComp
	}

	// Roll up stations
	stationsRollUp := make(map[string][]interchange.MergedMatch)
	for _, match := range matches {
		hashLol := fmt.Sprintf("%s%s%s%s", match.Competition, match.Date, match.Time, match.Title)
		_, prs := stationsRollUp[hashLol]
		if prs {
			stationsRollUp[hashLol] = append(stationsRollUp[hashLol], match)
		} else {
			stationsRollUp[hashLol] = []interchange.MergedMatch{match}
		}
	}
	matches = make([]interchange.MergedMatch, 0)
	for _, v := range stationsRollUp {
		if len(v) > 1 {
			stations := make([]string, 0)
			for _, foo := range v {
				stations = append(stations, foo.Station)
			}
			smoshed := v[0]
			smoshed.Station = strings.Join(stations, " | ")
			matches = append(matches, smoshed)
		} else {
			matches = append(matches, v[0])
		}
	}

	// Roll up dates
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

	for k, v := range matchesRollup {
		md := interchange.MergedMatchDay{Date: k, Matches: v}
		mergedFeed = append(mergedFeed, md)
	}

	sort.Slice(mergedFeed, func(i, j int) bool {
		return mergedFeed[i].Matches[0].Datetime < mergedFeed[j].Matches[0].Datetime
	})

	for _, matchDay := range mergedFeed {
		sort.Slice(matchDay.Matches, func(i, j int) bool {
			if matchDay.Matches[i].Time == matchDay.Matches[j].Time {
				return stationRank(matchDay.Matches[i].Station) < stationRank(matchDay.Matches[j].Station)
			}
			return matchDay.Matches[i].Time < matchDay.Matches[j].Time
		})
	}
	return mergedFeed
}

func shouldSkip(s string) bool {
	return strings.Contains(s, "Scottish")
}

func stationRank(station string) int {
	switch strings.Split(station, " | ")[0] {
	case "talkSPORT":
		return 1
	case "talkSPORT2":
		return 2
	case "BBC Radio 5":
		return 3
	default:
		return 4
	}
}

func mapTeamName(name string) string {
	nameMapper := map[string]string{
		"IR Iran":        "Iran",
		"Korea Republic": "South Korea",
		"Milan":          "AC Milan",
	}
	newName, prs := nameMapper[name]
	if prs {
		return newName
	} else {
		return name
	}
}

func mapCompName(comp string) string {
	comp = strings.TrimSuffix(comp, " Football 2022-23")

	return comp
}

func MergedMatchDayToEventList(mergedMatches []interchange.MergedMatchDay) []interchange.CalEvent {
	events := make([]interchange.CalEvent, 0)
	for _, day := range mergedMatches {
		for _, match := range day.Matches {
			starttime, err := time.Parse(time.RFC3339, match.Datetime)
			if err != nil {
				log.Fatalln(err)
			}
			event := interchange.CalEvent{
				Uid:      strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s/%s", match.Title, match.Competition)), " ", ""),
				DtStart:  starttime.UTC().Format(interchange.CalTimeString),
				Summary:  fmt.Sprintf("%s [%s]", match.Title, match.Competition),
				Location: match.Station,
			}
			events = append(events, event)
		}
	}
	return events
}
