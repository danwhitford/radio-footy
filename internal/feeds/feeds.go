package feeds

import (
	"log"

	"whitford.io/radiofooty/internal/urlgetter"
)

type matchGetter interface {
	getMatches() ([]Broadcast, error)
}

const (
	CalTimeString = "20060102T150405Z"
	niceDate      = "Monday, Jan 2"
	timeLayout    = "15:04"
)

func GetMatches() ([]MatchDay, error) {
	var broadcasts []Broadcast

	httpGetter := urlgetter.NewHttpGetter()
	macthGetters := []matchGetter{
		talkSportGetter{httpGetter},
		bbcMatchGetter{httpGetter},
		tvMatchGetter{httpGetter},
		f1MatchGetter{httpGetter},
		nflMatchGetter{httpGetter},
		manualGetter{},
	}

	for _, matchGetter := range macthGetters {
		got, err := matchGetter.getMatches()
		if err != nil {
			log.Println(err)
			continue
		}
		broadcasts = append(broadcasts, got...)
	}

	return matchesToMatchDays(broadcasts)
}

func matchesToMatchDays(broadcasts []Broadcast) ([]MatchDay, error) {
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
