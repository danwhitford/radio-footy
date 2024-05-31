package feeds

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"whitford.io/radiofooty/internal/broadcast"
)

func TestFilterMatches(t *testing.T) {
	table := []struct {
		input  []broadcast.Broadcast
		output []broadcast.Broadcast
	}{
		{
			input: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Chelsea",
						AwayTeam:    "Milan",
						Competition: "Premier League",
						Datetime: time.Date(2024, 6, 1, 15, 0, 0, 0, time.UTC),						
					},
				},
				{
					Match: broadcast.Match{
						HomeTeam:    "Inverness",
						AwayTeam:    "Hibernian",
						Competition: "Scottish Premiership",
						Datetime: time.Date(2024, 6, 1, 15, 0, 0, 0, time.UTC),						
					},
				},
				{
					Match: broadcast.Match{
						HomeTeam:    "Southampton",
						AwayTeam:    "Leeds",
						Competition: "Championship",
						Datetime: time.Date(2024, 5, 21, 16, 0, 0, 0, time.UTC),						
					},
				},
			},
			output: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Chelsea",
						AwayTeam:    "Milan",
						Competition: "Premier League",
						Datetime: time.Date(2024, 6, 1, 15, 0, 0, 0, time.UTC),						
					},
				},
			},
		},
	}

	fakeNow := time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC)
	for _, tst := range table {
		got := filterBroadcasts(tst.input, fakeNow)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestSortMatchDays(t *testing.T) {
	table := []struct {
		input  []broadcast.MatchDay
		output []broadcast.MatchDay
	}{
		{
			input: []broadcast.MatchDay{
				{
					DateTime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam: "Brentford",
								AwayTeam: "Arsenal",
								Datetime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam: "Bolton",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []broadcast.Station{broadcast.Talksport},
						},
						{
							Match: broadcast.Match{
								HomeTeam: "Fulham",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []broadcast.Station{broadcast.Talksport2},
						},
						{
							Match: broadcast.Match{
								HomeTeam: "Chelsea",
								AwayTeam: "Milan",
								Datetime: time.Date(2021, 8, 14, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam: "Romsey",
								AwayTeam: "Worthing",
								Datetime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
			output: []broadcast.MatchDay{
				{
					DateTime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam: "Chelsea",
								AwayTeam: "Milan",
								Datetime: time.Date(2021, 8, 14, 12, 0, 0, 0, time.UTC),
							},
						},
						{
							Match: broadcast.Match{
								HomeTeam: "Bolton",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []broadcast.Station{broadcast.Talksport},
						},
						{
							Match: broadcast.Match{
								HomeTeam: "Fulham",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []broadcast.Station{broadcast.Talksport2},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam: "Romsey",
								AwayTeam: "Worthing",
								Datetime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam: "Brentford",
								AwayTeam: "Arsenal",
								Datetime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
		},
	}

	for _, tst := range table {
		broadcast.SortMatchDays(tst.input)
		if diff := cmp.Diff(tst.output, tst.input); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestMatchesToMatchDays(t *testing.T) {
	table := []struct {
		input  []broadcast.Broadcast
		output []broadcast.MatchDay
	}{
		{
			input: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Chelsea",
						AwayTeam:    "Tottenham",
						Competition: "Premier League",
						Datetime:    time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Talksport,
				},
				{
					Match: broadcast.Match{
						HomeTeam:    "Chelsea",
						AwayTeam:    "Spurs",
						Competition: "Premier League",
						Datetime:    time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5,
				},
			},
			output: []broadcast.MatchDay{
				{
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam:    "Chelsea",
								AwayTeam:    "Tottenham",
								Competition: "Premier League",
								Datetime:    time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []broadcast.Station{broadcast.Talksport},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam:    "Chelsea",
								AwayTeam:    "Spurs",
								Competition: "Premier League",
								Datetime:    time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
							},
							Stations: []broadcast.Station{broadcast.Radio5},
						},
					},
				},
			},
		},
	}

	for _, tst := range table {
		got, err := matchesToMatchDays(tst.input)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestMatchDayToCalData(t *testing.T) {
	table := []struct {
		input  []broadcast.MatchDay
		output CalData
	}{
		{
			input: []broadcast.MatchDay{
				{
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []broadcast.Listing{
						{
							Match: broadcast.Match{
								HomeTeam:    "Chelsea",
								AwayTeam:    "Tottenham",
								Competition: "Premier League",
								Datetime:    time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []broadcast.Station{broadcast.Talksport, broadcast.Radio5},
						},
					},
				},
			},
			output: CalData{
				Events: []CalEvent{
					{
						Uid:      "chelseavtottenham/premierleague",
						DtStart:  "20210814T150000Z",
						Summary:  "Chelsea v Tottenham [Premier League]",
						Location: []broadcast.Station{broadcast.Talksport, broadcast.Radio5},
					},
				},
			},
		},
	}

	opts := cmpopts.IgnoreFields(CalData{}, "DtStamp")

	for _, tst := range table {
		got := MatchDayToCalData(tst.input)
		if diff := cmp.Diff(tst.output, got, opts); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
