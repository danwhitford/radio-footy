package feeds

import (
	"log"

	"whitford.io/radiofooty/internal/urlgetter"
)

type MatchGetter func(getter urlgetter.UrlGetter) ([]Broadcast, error)

const CalTimeString = "20060102T150405Z"
const niceDate = "Monday, Jan 2"
const timeLayout = "15:04"

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

	// Roll up stations
	listings := listingsFromBroadcasts(broadcasts)

	// Roll up dates
	mergedFeed := matchDaysFromListings(listings)

	// Sort by date, time, competition, title
	mergedFeed = sortMatchDays(mergedFeed)

	return mergedFeed
}

func filterMatches(matches []Broadcast) []Broadcast {
	filtered := make([]Broadcast, 0)
	for _, match := range matches {
		if match.Match.shouldSkip() {
			continue
		}
		filtered = append(filtered, match)
	}
	return filtered
}
