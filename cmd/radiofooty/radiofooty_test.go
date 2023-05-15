package main

import (
	"bytes"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/interchange"
)

func TestWriteIndex(t *testing.T) {
	table := []struct {
		input struct {
			MatchDays []interchange.MergedMatchDay
		}
		output string
	}{
		{
			input: struct {
				MatchDays []interchange.MergedMatchDay
			}{
				MatchDays: []interchange.MergedMatchDay{
					{
						NiceDate: "Monday, May 15",
						DateTime: time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
						Matches: []interchange.MergedMatch{
							{
								Title:       "Southampton v Manchester City",
								Datetime:    "2023-05-15T15:00:00Z",
								Competition: "Premier League",
								Station:     "BBC Radio 5 Live",
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
                    <p class="row text-row"><b>15:00 | BBC Radio 5 Live</b></p>
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
