package feeds

import (
	"log"

	"whitford.io/radiofooty/internal/urlgetter"
)

type matchGetter interface {
	getMatches() ([]Broadcast, error)
}

const (
	niceDate   = "Monday, Jan 2"
	timeLayout = "15:04"
)

func GetMatches() ([]MatchDay, error) {
	var broadcasts []Broadcast

	httpGetter := urlgetter.NewHttpGetter()
	macthGetters := []matchGetter{
		talkSportGetter{httpGetter},
		bbcMatchGetter{httpGetter},
		tvMatchGetter{httpGetter},
		skyGetter{httpGetter},
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

	days, err := matchesToMatchDays(broadcasts)
	if err != nil {
		return days, err
	}

	for _, d := range days {
		for _, similar := range d.reportSimilarGames(3) {
			log.Printf("'%v' is similar to '%v'\n", similar[0], similar[1])
		}
	}

	return days, nil
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
