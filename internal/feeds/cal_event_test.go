package feeds

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"whitford.io/radiofooty/internal/broadcast"
)

func TestUuidIsU(t *testing.T) {
	table := []struct{
		input []broadcast.MatchDay
		want CalData
	}{
		{
			input: []broadcast.MatchDay{
				{
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam:    "England",
								AwayTeam:    "West Indies",
								Competition: "Test Cricket",
								Datetime:    time.Date(2021, 8, 14, 8, 0, 0, 0, time.UTC),
							},
							Stations: []broadcast.Station{broadcast.Radio5},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam:    "England",
								AwayTeam:    "West Indies",
								Competition: "Test Cricket",
								Datetime:    time.Date(2021, 8, 15, 8, 0, 0, 0, time.UTC),
							},
							Stations: []broadcast.Station{broadcast.Radio5},
						},
					},
				},
			},
			want: CalData{
				Events: []CalEvent{
					{
						Uid:      "englandvwestindies/testcricket",
						DtStart:  "20210814T080000Z",
						Summary:  "England v West Indies [Test Cricket]",
						Location: []broadcast.Station{broadcast.Radio5},
					},
					{
						Uid:      "englandvwestindies/testcricket1",
						DtStart:  "20210815T080000Z",
						Summary:  "England v West Indies [Test Cricket]",
						Location: []broadcast.Station{broadcast.Radio5},
					},
				},
			},
		},
	}

	opts := cmpopts.IgnoreFields(CalData{}, "DtStamp")

	for _, tst := range table {
		got := MatchDayToCalData(tst.input)
		if diff := cmp.Diff(tst.want, got, opts); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}