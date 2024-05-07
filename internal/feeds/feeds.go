package feeds

import (
	"log"

	"whitford.io/radiofooty/internal/urlgetter"
	"whitford.io/radiofooty/internal/broadcast"
	"whitford.io/radiofooty/internal/channel"
)

type matchGetter interface {
	GetMatches() ([]broadcast.Broadcast, error)
}

const (
	niceDate   = "Monday, Jan 2"
	timeLayout = "15:04"
)

func GetMatches() ([]broadcast.MatchDay, error) {
	var broadcasts []broadcast.Broadcast

	httpGetter := urlgetter.NewHttpGetter()
	macthGetters := []matchGetter{
		channel.TalkSportGetter{Urlgetter: httpGetter},
		channel.BbcMatchGetter{Urlgetter: httpGetter},
		channel.TvMatchGetter{Urlgetter: httpGetter},
		channel.SkyGetter{Urlgetter: httpGetter},
		channel.ManualGetter{},
	}

	for _, matchGetter := range macthGetters {
		got, err := matchGetter.GetMatches()
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
		for _, similar := range d.ReportSimilarGames(3) {
			log.Printf("'%v' is similar to '%v'\n", similar[0], similar[1])
		}
	}

	return days, nil
}

func matchesToMatchDays(broadcasts []broadcast.Broadcast) ([]broadcast.MatchDay, error) {
	// Filter out matches we don't want
	broadcasts = filterBroadcasts(broadcasts)

	// Roll up stations
	listings := broadcast.ListingsFromBroadcasts(broadcasts)

	// Roll up dates
	mergedFeed, err := broadcast.MatchDaysFromListings(listings)
	if err != nil {
		return mergedFeed, err
	}

	// Sort by date, time, competition, title
	broadcast.SortMatchDays(mergedFeed)

	return mergedFeed, nil
}

func filterBroadcasts(broadcasts []broadcast.Broadcast) []broadcast.Broadcast {
	filtered := make([]broadcast.Broadcast, 0)
	for _, match := range broadcasts {
		if match.Match.ShouldSkip() {
			continue
		}
		filtered = append(filtered, match)
	}
	return filtered
}
