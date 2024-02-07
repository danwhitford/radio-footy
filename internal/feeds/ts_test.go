package feeds

import (
	"testing"
	"time"

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
						Datetime:    time.Date(2020, 12, 26, 17, 30, 0, 0, time.UTC),
						HomeTeam:    "Arsenal",
						AwayTeam:    "Chelsea",
						Competition: "Premier League",
					},
					Station: Talksport,
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
		{
			input: []TSGame{
				{
					HomeTeam: "Arsenal",
					AwayTeam: "West Bromwich Albion",
					League:   "Premier League Football 2023-24",
					Date:     "2024-02-07 20:00:00",
					Title:    "Arsenal v West Bromwich Albion",
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
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: Talksport,
				},
			},
		},
	}

	for _, test := range table {
		got := tsFeedToMatches(test.input)
		if diff := cmp.Diff(test.output, got); diff != "" {
			t.Errorf("tsFeedToMatches(%v) mismatch (-want +got):\n%s", test.input, diff)
		}

	}
}
