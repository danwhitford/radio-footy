package feeds

import (
	"log"
	"time"

	"whitford.io/radiofooty/internal/broadcast"
	"whitford.io/radiofooty/internal/channel"
	"whitford.io/radiofooty/internal/urlgetter"
)

type matchGetter interface {
	GetMatches() ([]broadcast.Broadcast, error)
}

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

	// Filter out matches we don't want
	broadcasts = filterBroadcasts(
		broadcasts,
		time.Date(
			time.Now().Year(),
			time.Now().Month(),
			time.Now().Day(),
			0,
			0,
			0,
			0,
			time.Now().Location(),
		),
	)

	days, err := matchesToMatchDays(broadcasts)
	if err != nil {
		return days, err
	}

	return days, nil
}

func matchesToMatchDays(broadcasts []broadcast.Broadcast) ([]broadcast.MatchDay, error) {
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

func filterBroadcasts(broadcasts []broadcast.Broadcast, now time.Time) []broadcast.Broadcast {
	filtered := make([]broadcast.Broadcast, 0)
	for _, match := range broadcasts {
		if match.Datetime.Before(now) {
			continue
		}
		if match.ShouldSkip() {
			continue
		}
		filtered = append(filtered, match)
	}
	return filtered
}
