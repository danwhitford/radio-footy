package feeds

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/interchange"
)

const niceDate = "Monday, Jan 2"
const timeLayout = "15:04"

type MatchPipeliner func(m interchange.MergedMatch) interchange.MergedMatch

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
		mapTeamNames(&matches[i])
		mapCompName(&matches[i])
	}

	// Roll up stations
	matches = rollUpStations(matches)

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

	// Sort by date
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
	case "BBC Radio 5":
		return 2
	case "talkSPORT2":
		return 3
	default:
		return 4
	}
}

func mapTeamNames(match *interchange.MergedMatch) {
	teams := strings.Split(match.Title, " v ")
	newTitle := fmt.Sprintf("%s v %s", mapTeamName(teams[0]), mapTeamName(teams[1]))
	match.Title = newTitle
}

func mapTeamName(name string) string {
	nameMapper := map[string]string{
		"IR Iran":                  "Iran",
		"Korea Republic":           "South Korea",
		"Milan":                    "AC Milan",
		"FC Bayern München":        "Bayern Munich",
		"Brighton and Hove Albion": "Brighton & Hove Albion",
		"Internazionale":           "Inter Milan",
	}
	newName, prs := nameMapper[name]
	if prs {
		return newName
	} else {
		return name
	}
}

func mapCompName(match *interchange.MergedMatch)  {
	match.Competition = strings.TrimSuffix(match.Competition, " Football 2022-23")
}

func rollUpStations(matches []interchange.MergedMatch) []interchange.MergedMatch {
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
			sort.Slice(stations, func(i, j int) bool {
				return stationRank(stations[i]) < stationRank(stations[j])
			})
			smoshed.Station = strings.Join(stations, " | ")
			matches = append(matches, smoshed)
		} else {
			matches = append(matches, v[0])
		}
	}
	return matches
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
