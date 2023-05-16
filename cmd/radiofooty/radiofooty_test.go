package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/feeds"
)

func TestWriteIndex(t *testing.T) {
	table := []struct {
		input struct {
			MatchDays []feeds.MergedMatchDay
		}
		output string
	}{
		{
			input: struct {
				MatchDays []feeds.MergedMatchDay
			}{
				MatchDays: []feeds.MergedMatchDay{
					{
						NiceDate: "Monday, May 15",
						DateTime: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
						Matches: []feeds.MergedMatch{
							{
								Title:       "Southampton v Manchester City",
								Datetime:    "2023-05-15T15:00:00Z",
								Competition: "Premier League",
								Stations:     []string{"talkSPORT", "BBC Radio 5 Live"},
								Time:        "15:00",
								Date:        "Monday, May 15",
							},
						},
					},
				},
			},
			output: `<!DOCTYPE html>
<html lang="en-gb">

<head>
    <title>Wireless Football</title>
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="Description" content="Upcoming football matches on the wireless">
    <meta content="text/html;charset=utf-8" http-equiv="Content-Type">
    <meta content="utf-8" http-equiv="encoding">
</head>

<body>
    <div id="container">
        <h1>Football on the radio</h1>
        
        <div class="matchday">
            <h2>Monday, May 15</h2>
            
                <div class="match">
                    <p class="row text-row"><b>15:00 | talkSPORT | BBC Radio 5 Live</b></p>
                    <p class="row text-row">Southampton v Manchester City (Premier League)</p>
                </div>
            
        </div>
        <hr />
        
    </div>
</body>
</html>
`,
		},
	}

	for _, tst := range table {
		var buffer bytes.Buffer
		writeIndex(tst.input, "../../internal/website/template.go.tmpl", &buffer)

		if diff := cmp.Diff(tst.output, buffer.String()); diff != "" {
			t.Errorf("writeIndex(%v) mismatch (-want +got):\n%s", tst.input, diff)
		}
	}
}

func TestWriteCal(t *testing.T) {
	table := []struct{
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
						Uid:       "leicestercityvliverpool/premierleague",
						DtStart:   "20230515T150000Z",
						Summary:  "Leicester City v Liverpool",
						Location: []string{"talkSPORT", "BBC Radio 5 Live"},
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
"LOCATION:talkSPORT | BBC Radio 5 Live\r\n" +
"DTSTAMP:20230515T205714Z\r\n" +
"DTSTART:20230515T150000Z\r\n" +
"DURATION:PT2H\r\n" +
"END:VEVENT\r\n" +
"END:VCALENDAR\r\n",
		},
	}

	for _, tst := range table {
		var buffer bytes.Buffer
		writeCal(tst.input, "../../internal/website/icalendar.go.tmpl", &buffer)

		if diff := cmp.Diff(tst.output, buffer.String()); diff != "" {
			t.Errorf("writeCal(%v) mismatch (-want +got):\n%s", tst.input, diff)
		}
	}
}