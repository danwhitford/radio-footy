package feeds

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"
)

const niceDate = "Monday, Jan 2"
const timeLayout = "15:04"

type MatchGetter func() ([]MergedMatch, error)

func GetMergedMatches() ([]MergedMatchDay, error) {
	var matches []MergedMatch
	getters := []MatchGetter{
		getTalkSportMatches,
		getBBCMatches,
		getSolentMatches,
		getTvMatches,
	}
	for _, getter := range getters {
		got, err := getter()
		if err != nil {
			return nil, err
		}
		matches = append(matches, got...)
	}

	return mergedMatchesToMergedMatchDays(matches), nil
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

func shouldSkip(m MergedMatch) bool {
	return strings.Contains(m.Competition, "Scottish") ||
		strings.Contains(m.Competition, "Women") ||
		strings.Contains(m.Title, "Scottish") ||
		strings.Contains(m.Title, "Women")
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
		return 99
	}
}

func mapTeamNames(match *MergedMatch) {
	teams := strings.Split(match.Title, " v ")
	if len(teams) != 2 {
		log.Fatalf("Got match with bad title: %+v", match)
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
	match.Competition = strings.TrimSuffix(match.Competition, " Football 2023-24")
	match.Competition = strings.TrimPrefix(match.Competition, "FA ")
	if match.Competition == "Test Match Special" {
		match.Competition = "The Ashes"
	}
}

func rollUpStations(matches []MergedMatch) []MergedMatch {
	stationsRollUp := make(map[string][]MergedMatch)
	for _, match := range matches {
		fmt.Printf("%+v\n", match)
		hashLol := fmt.Sprintf("%s%s%s", match.Competition, match.Date, match.Title)
		stationsRollUp[hashLol] = append(stationsRollUp[hashLol], match)
	}
	matches = make([]MergedMatch, 0)
	for _, v := range stationsRollUp {
		if len(v) > 1 {
			stations := make([]string, 0)
			events := make([]MergedMatchRadioEvent, 0)
			for _, foo := range v {
				stations = append(stations, foo.Stations...)
			}
			for _, foo := range v {
				event := MergedMatchRadioEvent{
					Station: foo.Stations[0],
					Date:    foo.Date,
					Time:    foo.Time,
				}
				events = append(events, event)
			}
			smoshed := v[0]
			sort.Slice(stations, func(i, j int) bool {
				return stationRank(stations[i]) < stationRank(stations[j])
			})
			sort.Slice(events, func(i, j int) bool {
				return stationRank(events[i].Station) < stationRank(events[j].Station)
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
		if shouldSkip(match) {
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

	for _, matches := range matchesRollup {
		if len(matches) == 1 {
			merged = append(merged, matches[0])
			continue
		}

		toCheck := make([]MergedMatch, 0)
		toCheck = append(toCheck, matches...)
		for len(toCheck) > 0 {
			candidate := toCheck[0]
			toCheck = toCheck[1:]
			matched := false
			for i, other := range toCheck {
				m1Teams := strings.Split(candidate.Title, " v ")
				m2Teams := strings.Split(other.Title, " v ")
				if m1Teams[0] == m2Teams[0] || m1Teams[1] == m2Teams[1] {
					stationList := make([]string, 0)
					stationList = append(stationList, candidate.Stations...)
					stationList = append(stationList, other.Stations...)
					sort.Slice(stationList, func(i, j int) bool {
						return stationRank(stationList[i]) < stationRank(stationList[j])
					})
					candidate.Stations = stationList
					merged = append(merged, candidate)
					toCheck = append(toCheck[:i], toCheck[i+1:]...)
					matched = true
					break
				}
			}
			if !matched {
				merged = append(merged, candidate)
			}
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
