package feeds

import (
	"log"
	"sort"
	"time"
)

type MatchDay struct {
	DateTime time.Time
	Matches  []Listing
}

func (m MatchDay) NiceDate() string {
	return m.DateTime.Format(niceDate)
}

func (m MatchDay) DateOnly() string {
	return m.DateTime.Format(time.DateOnly)
}

func matchDaysFromListings(matches []Listing) []MatchDay {
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

func sortMatchDays(matchDays []MatchDay) []MatchDay {
	sort.Slice(matchDays, func(i, j int) bool {
		return matchDays[i].DateTime.Before(matchDays[j].DateTime)
	})

	// Sort by time and station
	for _, matchDay := range matchDays {
		sort.Slice(matchDay.Matches, func(i, j int) bool {
			if matchDay.Matches[i].Datetime.Compare(matchDay.Matches[j].Datetime) == 0 {
				return matchDay.Matches[i].Rank() < matchDay.Matches[j].Rank()
			}
			return matchDay.Matches[i].Datetime.Before(matchDay.Matches[j].Datetime)
		})
	}

	return matchDays
}
