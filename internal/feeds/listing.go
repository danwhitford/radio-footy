package feeds

import (
	"fmt"
	"regexp"
	"sort"
	"time"
)

type Listing struct {
	Match
	Stations []Station
}

func (l Listing) GameHash() string {
	s := fmt.Sprintf("%s%s", l.Datetime.Format(time.RFC3339), l.Title())
	r := regexp.MustCompile("[^0-9a-zA-Z]")
	s = r.ReplaceAllString(s, "")
	return s
}

func listingsFromBroadcasts(broadcasts []Broadcast) []Listing {
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

func (l Listing) Rank() int {
	return l.Stations[0].Rank
}
