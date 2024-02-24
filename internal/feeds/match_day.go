package feeds

import (
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

func matchDaysFromListings(matches []Listing) ([]MatchDay, error) {
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
			return matchDays, err
		}
		md := MatchDay{Matches: matches, DateTime: dt}
		matchDays = append(matchDays, md)
	}

	return matchDays, nil
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

func (matchDay MatchDay) reportSimilarGames(threshold int) [][2]Listing {
	similars := make([][2]Listing, 0)
	for i, m := range matchDay.Matches {
		for _, other := range matchDay.Matches[i+1:] {
			score := m.similarityScore(other.Match)
			if score >= 2 {
				similars = append(similars, [2]Listing{
					m,
					other,
				})
			}
		}
	}
	return similars
}
