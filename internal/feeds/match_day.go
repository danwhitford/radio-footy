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

func sortMatchDays(matchDays []MatchDay) {
	sort.Slice(matchDays, func(i, j int) bool {
		return matchDays[i].DateTime.Before(matchDays[j].DateTime)
	})

	for _, matchDay := range matchDays {
		matchDay.sortMatchDay()
	}
}

func (matchDay *MatchDay) sortMatchDay() {
	sort.Slice(matchDay.Matches, func(i, j int) bool {
		return matchDay.Matches[i].less(matchDay.Matches[j])
	})
}
