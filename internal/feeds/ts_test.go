package feeds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTsFeedToMatches(t *testing.T) {
	table := []struct {
		input  []TSGame
		output []Broadcast
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
			output: []Broadcast{
				{
					Match: Match{
						Time:        "17:30",
						Date:        "Saturday, Dec 26",
						Datetime:    "2020-12-26T17:30:00Z",
						HomeTeam:    "Arsenal",
						AwayTeam:    "Chelsea",
						Competition: "Premier League",
					},
					Station: "talkSPORT",
				},
			},
		},
		{
			input: []TSGame{
				{
					Livefeed: []TSLiveFeed{
						{
							Feedname: "talkSPORT",
						},
					},
					Sport:    "Football",
					Date:     "2023-07-01 15:15:00",
					HomeTeam: "England",
					AwayTeam: "Portugal",
					League:   "International Friendlies",
					Title:    "England Women v Portugal Women",
				},
			},
			output: []Broadcast{},
		},
	}

	for _, test := range table {
		got := tsFeedToMatches(test.input)
		if diff := cmp.Diff(test.output, got); diff != "" {
			t.Errorf("tsFeedToMatches(%v) mismatch (-want +got):\n%s", test.input, diff)
		}

	}
}
