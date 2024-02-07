package feeds

import (
	"fmt"
	"strings"
)

type CalEvent struct {
	Uid      string
	DtStart  string
	Summary  string
	Location []Station
}

func MatchDayToEventList(Matches []MatchDay) []CalEvent {
	events := make([]CalEvent, 0)
	for _, day := range Matches {
		for _, match := range day.Matches {
			starttime := match.Datetime

			event := CalEvent{
				Uid:      strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s/%s", match.Title(), match.Competition)), " ", ""),
				DtStart:  starttime.UTC().Format(CalTimeString),
				Summary:  fmt.Sprintf("%s [%s]", match.Title(), match.Competition),
				Location: match.Stations,
			}
			events = append(events, event)
		}
	}
	return events
}
