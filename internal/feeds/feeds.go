package feeds

import (
	"log"

	"whitford.io/radiofooty/internal/urlgetter"
)

type MatchGetter func(getter urlgetter.UrlGetter) ([]Broadcast, error)

const (
	CalTimeString = "20060102T150405Z"
	niceDate      = "Monday, Jan 2"
	timeLayout    = "15:04"
)

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

	return MatchesToMatchDays(broadcasts)
}

func MatchesToMatchDays(broadcasts []Broadcast) ([]MatchDay, error) {
	// Filter out matches we don't want
	broadcasts = filterBroadcasts(broadcasts)

	// Roll up stations
	listings := listingsFromBroadcasts(broadcasts)

	// Roll up dates
	mergedFeed, err := matchDaysFromListings(listings)
	if err != nil {
		return mergedFeed, err
	}

	// Sort by date, time, competition, title
	sortMatchDays(mergedFeed)

	return mergedFeed, nil
}

func filterBroadcasts(broadcasts []Broadcast) []Broadcast {
	filtered := make([]Broadcast, 0)
	for _, match := range broadcasts {
		if match.Match.shouldSkip() {
			continue
		}
		filtered = append(filtered, match)
	}
	return filtered
}
