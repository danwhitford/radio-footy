package feeds

import (
	"fmt"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/broadcast"
)

type CalEvent struct {
	Uid      string
	DtStart  string
	Summary  string
	Location []broadcast.Station
}

type CalData struct {
	Events  []CalEvent
	DtStamp string
}

const calTimeString string = "20060102T150405Z"

func MatchDayToCalData(Matches []broadcast.MatchDay) CalData {
	uuidCount := make(map[string]int)

	events := make([]CalEvent, 0)
	for _, day := range Matches {
		for _, match := range day.Matches {
			starttime := match.Datetime
			uidRoot := strings.ReplaceAll(strings.ToLower(fmt.Sprintf("%s/%s", match.Title(), match.Competition)), " ", "")
			var uid string

			if _, prs := uuidCount[uidRoot]; prs {
				uid = fmt.Sprintf("%s%d", uidRoot, uuidCount[uidRoot])
				uuidCount[uidRoot]++
			} else {
				uid = uidRoot
				uuidCount[uidRoot] = 1
			}

			event := CalEvent{
				Uid:      uid,
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
