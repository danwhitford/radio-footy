package channel

import (
	"testing"
	"time"
	"whitford.io/radiofooty/internal/broadcast"

	"github.com/google/go-cmp/cmp"
)

func TestTsFeedToMatches(t *testing.T) {
	table := []struct {
		input  []TSGame
		output []broadcast.Broadcast
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
			output: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						Datetime:    time.Date(2020, 12, 26, 17, 30, 0, 0, time.UTC),
						HomeTeam:    "Arsenal",
						AwayTeam:    "Chelsea",
						Competition: "Premier League",
					},
					Station: broadcast.Talksport,
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
			output: []broadcast.Broadcast{},
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
			output: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Talksport,
				},
			},
		},
		{
			input: []TSGame{
				{
					Livefeed: []TSLiveFeed{
						{
							Feedname: "talkSPORT2",
						}},
					Sport:    "Football",
					Date:     "2024-02-27 17:00:00",
					HomeTeam: "England",
					AwayTeam: "Italy",
					League:   "International Friendlies", // Not international men's football
					Title:    "Friendly - England v Italy",
				},
			},
			output: []broadcast.Broadcast{},
		},
	}

	for _, test := range table {
		got := tsFeedToMatches(test.input)
		if diff := cmp.Diff(test.output, got); diff != "" {
			t.Errorf("tsFeedToMatches(%v) mismatch (-want +got):\n%s", test.input, diff)
		}
	}
}
