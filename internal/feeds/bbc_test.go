package feeds

import (
	"testing"
	"whitford.io/radiofooty/internal/interchange"
	"github.com/google/go-cmp/cmp"
)

func TestBBCDayToMergedMatch(t *testing.T) {
	table := []struct {
		input  interchange.BBCFeed
		output []interchange.MergedMatch
	}{
		{
			input: interchange.BBCFeed{
				Data: []interchange.BBCFeedData{
					{
						Data: []interchange.BBCProgramData{
							{
								Title: interchange.BBCTitles{
									Primary:   "5 Live Sport",
									Secondary: "Premier League Football 2022-23",
									Tertiary:  "Arsenal v Chelsea",
								},
								Start: "2020-12-26T17:30:00Z",
							},
						},
					},
				},
			},
			output: []interchange.MergedMatch{
				{
					Title: 	"Arsenal v Chelsea",
					Stations: []string{"BBC Radio 5"},
					Competition: "Premier League Football 2022-23",
					Time:     "17:30",
					Date:     "Saturday, Dec 26",
					Datetime: "2020-12-26T17:30:00Z",
				},
			},
		},
	}

	for _, test := range table {
		got := bbcDayToMergedMatch(test.input)
		if diff := cmp.Diff(test.output, got); diff != "" {
			t.Errorf("bbcDayToMergedMatch() mismatch (-want +got):\n%s", diff)
		}
	}
}