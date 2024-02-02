package main

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/feeds"
)

func TestWriteCal(t *testing.T) {
	table := []struct {
		input struct {
			DtStamp string
			Events  []feeds.CalEvent
		}
		output string
	}{
		{
			input: struct {
				DtStamp string
				Events  []feeds.CalEvent
			}{
				DtStamp: "20230515T205714Z",
				Events: []feeds.CalEvent{
					{
						Uid:      "leicestercityvliverpool/premierleague",
						DtStart:  "20230515T150000Z",
						Summary:  "Leicester City v Liverpool",
						Location: []feeds.Station{feeds.Talksport, feeds.Radio5},
					},
				},
			},
			output: "BEGIN:VCALENDAR\r\n" +
				"VERSION:2.0\r\n" +
				"METHOD:PUBLISH\r\n" +
				"PRODID:-wirelessfootball.co.uk/icalendar\r\n" +
				"BEGIN:VEVENT\r\n" +
				"UID:leicestercityvliverpool/premierleague\r\n" +
				"SUMMARY:Leicester City v Liverpool\r\n" +
				"DESCRIPTION:Leicester City v Liverpool\r\n" +
				"LOCATION:talkSPORT | Radio 5 Live\r\n" +
				"DTSTAMP:20230515T205714Z\r\n" +
				"DTSTART:20230515T150000Z\r\n" +
				"DURATION:PT2H\r\n" +
				"END:VEVENT\r\n" +
				"END:VCALENDAR\r\n",
		},
	}

	for _, tst := range table {
		var buffer bytes.Buffer
		writeIndex(tst.input, "icalendar.go.tmpl", "../../internal/website/icalendar.go.tmpl", &buffer)

		if diff := cmp.Diff(tst.output, buffer.String()); diff != "" {
			t.Errorf("writeCal(%v) mismatch (-want +got):\n%s", tst.input, diff)
		}
	}
}
