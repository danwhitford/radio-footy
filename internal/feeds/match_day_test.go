package feeds

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMatchDaysFromListings(t *testing.T) {
	table := []struct {
		input  []Listing
		output []MatchDay
	}{
		{
			input: []Listing{
				{
					Match: Match{
						HomeTeam: "Milan",
						AwayTeam: "Chelsea",
						Datetime: time.Date(2021, 8, 14, 19, 0, 0, 0, time.UTC),
					},
				},
				{
					Match: Match{
						HomeTeam: "Bolton",
						AwayTeam: "Barnsley",
						Datetime: time.Date(2021, 8, 14, 21, 0, 0, 0, time.UTC),
					},
				},
				{
					Match: Match{
						HomeTeam: "Romsey",
						AwayTeam: "Worthing",
						Datetime: time.Date(2021, 8, 15, 19, 0, 0, 0, time.UTC),
					},
				},
			},
			output: []MatchDay{
				{
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								Datetime: time.Date(2021, 8, 14, 19, 0, 0, 0, time.UTC),
								HomeTeam: "Milan",
								AwayTeam: "Chelsea",
							},
						},
						{
							Match: Match{
								Datetime: time.Date(2021, 8, 14, 21, 0, 0, 0, time.UTC),
								HomeTeam: "Bolton",
								AwayTeam: "Barnsley",
							},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								Datetime: time.Date(2021, 8, 15, 19, 0, 0, 0, time.UTC),
								HomeTeam: "Romsey",
								AwayTeam: "Worthing",
							},
						},
					},
				},
			},
		},
	}

	for _, tst := range table {
		less := func(i, j MatchDay) bool {
			return i.DateTime.Before(j.DateTime)
		}
		got, err := matchDaysFromListings(tst.input)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(tst.output, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
