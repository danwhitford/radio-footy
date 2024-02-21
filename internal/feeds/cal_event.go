package feeds

import (
	"fmt"
	"strings"
	"time"
)

type CalEvent struct {
	Uid      string
	DtStart  string
	Summary  string
	Location []Station
}

type CalData struct {
	Events  []CalEvent
	DtStamp string
}

const calTimeString string = "20060102T150405Z"

func MatchDayToCalData(Matches []MatchDay) CalData {
	events := make([]CalEvent, 0)
	for _, day := range Matches {
		for _, match := range day.Matches {
			starttime := match.Datetime

			event := CalEvent{
				Uid:      strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s/%s", match.Title(), match.Competition)), " ", ""),
				DtStart:  starttime.UTC().Format(calTimeString),
				Summary:  fmt.Sprintf("%s [%s]", match.Title(), match.Competition),
				Location: match.Stations,
			}
			events = append(events, event)
		}
	}
	return CalData{
		Events:  events,
		DtStamp: time.Now().Format(calTimeString),
	}
}
