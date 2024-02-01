package feeds

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/urlgetter"
)

type MatchGetter func(getter urlgetter.UrlGetter) ([]Broadcast, error)

func GetMatches() ([]MatchDay, error) {
	var matches []Broadcast
	macthGetters := []MatchGetter{
		getTalkSportMatches,
		getBBCMatches,
		getTvMatches,
		getNflOnSky,
	}

	httpGetter := urlgetter.NewHttpGetter()
	for _, matchGetter := range macthGetters {
		got, err := matchGetter(httpGetter)
		if err != nil {
			return nil, err
		}
		matches = append(matches, got...)
	}

	return MatchesToMatchDays(matches), nil
}

func MatchesToMatchDays(matches []Broadcast) []MatchDay {
	// Filter out matches we don't want
	matches = filterMatches(matches)

	// Map team names and competition names
	for i := range matches {
		mapTeamNames(&matches[i].Match)
		mapCompName(&matches[i].Match)
	}

	// Roll up stations
	listings := rollUpStations(matches)

	// Roll up dates
	mergedFeed := rollUpDates(listings)

	// Sort by date, time, competition, title
	mergedFeed = sortMatchDays(mergedFeed)

	return mergedFeed
}

func shouldSkip(m Match) bool {
	return strings.Contains(m.Competition, "Scottish") ||
		strings.Contains(m.Competition, "Women") ||
		strings.Contains(m.HomeTeam, "Scottish") ||
		strings.Contains(m.HomeTeam, "Women")
}

func stationRank(station string) int {
	stationsInOrder := []string{
		"Sky Sports",
		"talkSPORT",
		"BBC Radio 5",
		"talkSPORT2",
	}
	for i, s := range stationsInOrder {
		if station == s {
			return i
		}
	}

	return 99
}

func mapTeamNames(match *Match) {
	match.HomeTeam = mapTeamName(match.HomeTeam)
	match.AwayTeam = mapTeamName(match.AwayTeam)
}

func mapTeamName(name string) string {
	nameMapper := map[string]string{
		"IR Iran":                  "Iran",
		"Korea Republic":           "South Korea",
		"Milan":                    "AC Milan",
		"FC Bayern München":        "Bayern Munich",
		"Brighton and Hove Albion": "Brighton & Hove Albion",
		"Internazionale":           "Inter Milan",
		"Wolverhampton Wanderers":  "Wolves",
		"West Bromwich Albion":     "West Brom",
	}
	newName, prs := nameMapper[name]
	if prs {
		return newName
	} else {
		return name
	}
}

func mapCompName(match *Match) {
	keepPrefixes := []string{
		"Premier League",
		"FA Cup",
		"UEFA Champions League",
	}
	for _, prefix := range keepPrefixes {
		if strings.HasPrefix(match.Competition, prefix) {
			match.Competition = prefix
			return
		}
	}

	replacements := map[string]string {
		"Carabao Cup": "EFL Cup",
		"English Football League Trophy": "EFL Cup",
	}
	for old, new := range replacements { 
		if match.Competition == old {
			match.Competition = new
			return
		}
	}
}

func rollUpStations(matches []Broadcast) []Listing {
	stationsRollUp := make(map[string][]Broadcast)
	for _, match := range matches {
		hashLol := fmt.Sprintf("%s%v%s%s", match.Competition, match.Datetime, match.HomeTeam, match.AwayTeam)
		stationsRollUp[hashLol] = append(stationsRollUp[hashLol], match)
	}
	listings := make([]Listing, 0)
	for _, v := range stationsRollUp {
		stations := make([]string, 0)
		for _, bcst := range v {
			stations = append(stations, bcst.Station)
		}
		sort.Slice(stations, func(i, j int) bool {
			return stationRank(stations[i]) < stationRank(stations[j])
		})
		listings = append(listings, Listing{
			Match:    v[0].Match,
			Stations: stations,
		})
	}
	return listings
}

func rollUpDates(matches []Listing) []MatchDay {
	matchesRollup := make(map[string][]Listing)
	for _, match := range matches {
		d := match.Datetime
		
		key := d.Format(time.DateOnly)
		matchesRollup[key] = append(matchesRollup[key], match)
	}

	matchDays := make([]MatchDay, 0)
	for k, matches := range matchesRollup {
		dt, err := time.Parse(time.DateOnly, k)
		if err != nil {
			log.Fatal(err)
		}
		md := MatchDay{Matches: matches, DateTime: dt}
		matchDays = append(matchDays, md)
	}

	return matchDays
}

func filterMatches(matches []Broadcast) []Broadcast {
	filtered := make([]Broadcast, 0)
	for _, match := range matches {
		if shouldSkip(match.Match) {
			continue
		}
		filtered = append(filtered, match)
	}
	return filtered
}

func sortMatchDays(matchDays []MatchDay) []MatchDay {
	sort.Slice(matchDays, func(i, j int) bool {
		return matchDays[i].DateTime.Before(matchDays[j].DateTime)
	})

	// Sort by time and station
	for _, matchDay := range matchDays {
		sort.Slice(matchDay.Matches, func(i, j int) bool {
			if matchDay.Matches[i].Datetime.Compare(matchDay.Matches[j].Datetime) == 0 {
				return stationRank(matchDay.Matches[i].Stations[0]) < stationRank(matchDay.Matches[j].Stations[0])
			}
			return matchDay.Matches[i].Datetime.Before(matchDay.Matches[j].Datetime)
		})
	}

	return matchDays
}

func MatchDayToEventList(Matches []MatchDay) []CalEvent {
	events := make([]CalEvent, 0)
	for _, day := range Matches {
		for _, match := range day.Matches {
			starttime := match.Datetime
			
			event := CalEvent{
				Uid:      strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s/%s", match.Title(), match.Competition)), " ", ""),
				DtStart:  starttime.UTC().Format(CalTimeString),
				Summary:  fmt.Sprintf("%s [%s]", match.Title(), match.Competition),
				Location: match.Stations,
			}
			events = append(events, event)
		}
	}
	return events
}
