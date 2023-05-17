package feeds

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"
)

const niceDate = "Monday, Jan 2"
const timeLayout = "15:04"

func GetMergedMatches() []MergedMatchDay {
	var matches []MergedMatch
	matches = append(matches, getTalkSportMatches()...)
	matches = append(matches, getBBCMatches()...)

	return mergedMatchesToMergedMatchDays(matches)
}

func mergedMatchesToMergedMatchDays(matches []MergedMatch) []MergedMatchDay {
	// Filter out matches we don't want
	matches = filterMatches(matches)

	// Map team names and competition names
	for i := range matches {
		mapTeamNames(&matches[i])
		mapCompName(&matches[i])
	}
	matches = fuzzyMergeTeams(matches)

	// Roll up stations
	matches = rollUpStations(matches)

	// Roll up dates
	mergedFeed := rollUpDates(matches)

	// Sort by date, time, competition, title
	mergedFeed = sortMatchDays(mergedFeed)

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

func mapTeamNames(match *MergedMatch) {
	teams := strings.Split(match.Title, " v ")
	if len(teams) != 2 {
		log.Printf("Skipping match with bad title: %v", match)
		os.Exit(1)
	}
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

func mapCompName(match *MergedMatch) {
	match.Competition = strings.TrimSuffix(match.Competition, " Football 2022-23")
}

func rollUpStations(matches []MergedMatch) []MergedMatch {
	stationsRollUp := make(map[string][]MergedMatch)
	for _, match := range matches {
		hashLol := fmt.Sprintf("%s%s%s%s", match.Competition, match.Date, match.Time, match.Title)
		_, prs := stationsRollUp[hashLol]
		if prs {
			stationsRollUp[hashLol] = append(stationsRollUp[hashLol], match)
		} else {
			stationsRollUp[hashLol] = []MergedMatch{match}
		}
	}
	matches = make([]MergedMatch, 0)
	for _, v := range stationsRollUp {
		if len(v) > 1 {
			stations := make([]string, 0)
			for _, foo := range v {
				stations = append(stations, foo.Stations...)
			}
			smoshed := v[0]
			sort.Slice(stations, func(i, j int) bool {
				return stationRank(stations[i]) < stationRank(stations[j])
			})
			smoshed.Stations = stations
			matches = append(matches, smoshed)
		} else {
			matches = append(matches, v[0])
		}
	}
	return matches
}

func rollUpDates(matches []MergedMatch) []MergedMatchDay {
	matchesRollup := make(map[string][]MergedMatch)
	for _, match := range matches {
		d, err := time.Parse(time.RFC3339, match.Datetime)
		if err != nil {
			log.Fatalf("error rolling up dates %s", err)
		}
		key := d.Format(time.DateOnly)
		matchesRollup[key] = append(matchesRollup[key], match)
	}

	matchDays := make([]MergedMatchDay, 0)
	for k, matches := range matchesRollup {
		dt, err := time.Parse(time.DateOnly, k)
		if err != nil {
			log.Fatal(err)
		}
		md := MergedMatchDay{NiceDate: dt.Format(niceDate), Matches: matches, DateTime: dt}
		matchDays = append(matchDays, md)
	}

	return matchDays
}

func filterMatches(matches []MergedMatch) []MergedMatch {
	filtered := make([]MergedMatch, 0)
	for _, match := range matches {
		if shouldSkip(match.Competition) {
			continue
		}
		filtered = append(filtered, match)
	}
	return filtered
}

func sortMatchDays(matchDays []MergedMatchDay) []MergedMatchDay {
	sort.Slice(matchDays, func(i, j int) bool {
		return matchDays[i].DateTime.Before(matchDays[j].DateTime)
	})

	// Sort by time and station
	for _, matchDay := range matchDays {
		sort.Slice(matchDay.Matches, func(i, j int) bool {
			if matchDay.Matches[i].Time == matchDay.Matches[j].Time {
				return stationRank(matchDay.Matches[i].Stations[0]) < stationRank(matchDay.Matches[j].Stations[0])
			}
			return matchDay.Matches[i].Time < matchDay.Matches[j].Time
		})
	}

	return matchDays
}

func fuzzyMergeTeams(matches []MergedMatch) []MergedMatch {
	merged := make([]MergedMatch, 0)
	matchesRollup := make(map[string][]MergedMatch)
	for _, match := range matches {
		key := fmt.Sprintf("%s%s", match.Competition, match.Datetime)
		matchesRollup[key] = append(matchesRollup[key], match)
	}

	for key, match := range matchesRollup {
		if len(match) > 2 {
			log.Fatalf("Too many matches for key %s: %v", key, match)
		}
		if len(match) == 1 {
			merged = append(merged, match[0])
			continue
		}

		m1Teams := strings.Split(match[0].Title, " v ")
		m2Teams := strings.Split(match[1].Title, " v ")
		if m1Teams[0] == m2Teams[0] || m1Teams[1] == m2Teams[1] {
			stationList := make([]string, 0)
			stationList = append(stationList, match[0].Stations...)
			stationList = append(stationList, match[1].Stations...)
			sort.Slice(stationList, func(i, j int) bool {
				return stationRank(stationList[i]) < stationRank(stationList[j])
			})
			match[0].Stations = stationList
			merged = append(merged, match[0])
		} else {
			merged = append(merged, match...)
		}
	}
	return merged
}

func MergedMatchDayToEventList(mergedMatches []MergedMatchDay) []CalEvent {
	events := make([]CalEvent, 0)
	for _, day := range mergedMatches {
		for _, match := range day.Matches {
			starttime, err := time.Parse(time.RFC3339, match.Datetime)
			if err != nil {
				log.Fatalln("error while creating event list", err)
			}
			event := CalEvent{
				Uid:      strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s/%s", match.Title, match.Competition)), " ", ""),
				DtStart:  starttime.UTC().Format(CalTimeString),
				Summary:  fmt.Sprintf("%s [%s]", match.Title, match.Competition),
				Location: match.Stations,
			}
			events = append(events, event)
		}
	}
	return events
}
