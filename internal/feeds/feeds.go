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
	var broadcasts []Broadcast
	macthGetters := []MatchGetter{
		getTalkSportMatches,
		getBBCMatches,
		getTvMatches,
		getNflOnSky,
		getManualFeeds,
	}

	httpGetter := urlgetter.NewHttpGetter()
	for _, matchGetter := range macthGetters {
		got, err := matchGetter(httpGetter)
		if err != nil {
			log.Println(err)
			continue
		}
		broadcasts = append(broadcasts, got...)
	}

	return MatchesToMatchDays(broadcasts), nil
}

func MatchesToMatchDays(broadcasts []Broadcast) []MatchDay {
	// Filter out matches we don't want
	broadcasts = filterMatches(broadcasts)

	// Map team names and competition names
	for i := range broadcasts {
		broadcasts[i].Match.mapTeamNames()
		broadcasts[i].Match.mapCompName()
	}

	// Roll up stations
	listings := rollUpStations(broadcasts)

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

func rollUpStations(broadcasts []Broadcast) []Listing {
	stationsRollUp := make(map[string]Listing)
	for _, bcst := range broadcasts {
		hashLol := bcst.RollUpHash()
		if listing, prs := stationsRollUp[hashLol]; prs {
			listing.Stations = append(listing.Stations, bcst.Station)
			if bcst.Datetime.After(listing.Datetime) {
				listing.Datetime = bcst.Datetime
			}
			stationsRollUp[hashLol] = listing
		} else {
			stationsRollUp[hashLol] = Listing{
				bcst.Match,
				[]Station{bcst.Station},
			}
		}
	}

	listings := make([]Listing, 0)
	for _, listing := range stationsRollUp {
		sort.Slice(listing.Stations, func(i, j int) bool {
			return listing.Stations[i].Rank < listing.Stations[j].Rank
		})
		listings = append(listings, listing)
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
				return matchDay.Matches[i].Stations[0].Rank < matchDay.Matches[j].Stations[0].Rank
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
