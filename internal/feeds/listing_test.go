package feeds

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestRollUpStations(t *testing.T) {
	table := []struct {
		input  []Broadcast
		output []Listing
	}{
		{
			input: []Broadcast{
				{
					Station: Radio5,
					Match:   Match{},
				},
				{
					Station: Talksport,
					Match:   Match{},
				},
			},
			output: []Listing{
				{
					Match:    Match{},
					Stations: []Station{Talksport, Radio5},
				},
			},
		},
		{
			input: []Broadcast{
				{
					Station: Radio5,
					Match:   Match{},
				},
				{
					Station: Radio5Extra,
					Match:   Match{},
				},
			},
			output: []Listing{
				{
					Match:    Match{},
					Stations: []Station{Radio5, Radio5Extra},
				},
			},
		},
		{
			input: []Broadcast{
				{
					Station: Talksport,
					Match:   Match{},
				},
				{
					Station: Radio5,
					Match:   Match{},
				},
				{
					Station: ChannelFour,
					Match:   Match{},
				},
			},
			output: []Listing{
				{
					Match:    Match{},
					Stations: []Station{ChannelFour, Talksport, Radio5},
				},
			},
		},
		{
			input: []Broadcast{
				{
					Station: Radio5Extra,
					Match: Match{
						HomeTeam: "West Ham United",
						AwayTeam: "AFC Bournemouth",
						Datetime: time.Date(2024, 2, 1, 19, 25, 0, 0, time.UTC),
					},
				},
				{
					Station: TNTSports,
					Match: Match{
						HomeTeam: "West Ham United",
						AwayTeam: "AFC Bournemouth",
						Datetime: time.Date(2024, 2, 1, 19, 30, 0, 0, time.UTC),
					},
				},
			},
			output: []Listing{
				{
					Stations: []Station{TNTSports, Radio5Extra},
					Match: Match{
						HomeTeam: "West Ham United",
						AwayTeam: "AFC Bournemouth",
						Datetime: time.Date(2024, 2, 1, 19, 30, 0, 0, time.UTC),
					},
				},
			},
		},
	}

	for _, tst := range table {
		got := listingsFromBroadcasts(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	}
}
