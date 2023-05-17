package feeds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTsFeedToMergedMatches(t *testing.T) {
	table := []struct {
		input  []TSGame
		output []MergedMatch
	}{
		{
			input: []TSGame{
				{
					HomeTeam: "Arsenal",
					AwayTeam: "Chelsea",
					League:   "Premier League",
					Date:     "2020-12-26 17:30:00",
					Title:    "Arsenal v Chelsea",
					Livefeed: []TSLiveFeed{
						{
							Feedname: "talkSPORT",
						},
					},
				},
			},
			output: []MergedMatch{
				{
					Time:        "17:30",
					Date:        "Saturday, Dec 26",
					Stations:    []string{"talkSPORT"},
					Datetime:    "2020-12-26T17:30:00Z",
					Title:       "Arsenal v Chelsea",
					Competition: "Premier League",
				},
			},
		},
	}

	for _, test := range table {
		got := tsFeedToMergedMatches(test.input)
		if diff := cmp.Diff(test.output, got); diff != "" {
			t.Errorf("tsFeedToMergedMatches(%v) mismatch (-want +got):\n%s", test.input, diff)
		}

	}
}
