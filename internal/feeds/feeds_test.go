package feeds

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"

)

func TestFilterMatches(t *testing.T) {
	table := []struct {
		input  []Broadcast
		output []Broadcast
	}{
		{
			input: []Broadcast{
				{
					Match: Match{
						HomeTeam:    "Chelsea",
						AwayTeam:    "Milan",
						Competition: "Premier League",
					},
				},
				{
					Match: Match{
						HomeTeam:    "Inverness",
						AwayTeam:    "Hibernian",
						Competition: "Scottish Premiership",
					},
				},
			},
			output: []Broadcast{
				{
					Match: Match{
						HomeTeam:    "Chelsea",
						AwayTeam:    "Milan",
						Competition: "Premier League",
					},
				},
			},
		},
	}

	for _, tst := range table {
		got := filterBroadcasts(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestSortMatchDays(t *testing.T) {
	table := []struct {
		input  []MatchDay
		output []MatchDay
	}{
		{
			input: []MatchDay{
				{
					DateTime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								HomeTeam: "Brentford",
								AwayTeam: "Arsenal",
								Datetime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								HomeTeam: "Bolton",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []Station{Talksport},
						},
						{
							Match: Match{
								HomeTeam: "Fulham",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []Station{Talksport2},
						},
						{
							Match: Match{
								HomeTeam: "Chelsea",
								AwayTeam: "Milan",
								Datetime: time.Date(2021, 8, 14, 12, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								HomeTeam: "Romsey",
								AwayTeam: "Worthing",
								Datetime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
							},
						},
					},
				},
			},
			output: []MatchDay{
				{
					DateTime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								HomeTeam: "Chelsea",
								AwayTeam: "Milan",
								Datetime: time.Date(2021, 8, 14, 12, 0, 0, 0, time.UTC),
							},
						},
						{
							Match: Match{
								HomeTeam: "Bolton",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []Station{Talksport},
						},
						{
							Match: Match{
								HomeTeam: "Fulham",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []Station{Talksport2},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								HomeTeam: "Romsey",
								AwayTeam: "Worthing",
								Datetime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
							},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
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
		sortMatchDays(tst.input)
		if diff := cmp.Diff(tst.output, tst.input); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestMatchesToMatchDays(t *testing.T) {
	table := []struct {
		input  []Broadcast
		output []MatchDay
	}{
		{
			input: []Broadcast{
				{
					Match: Match{
						HomeTeam:    "Chelsea",
						AwayTeam:    "Tottenham",
						Competition: "Premier League",
						Datetime:    time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					},
					Station: Talksport,
				},
				{
					Match: Match{
						HomeTeam:    "Inverness",
						AwayTeam:    "Hibernian",
						Competition: "Scottish Premiership",
						Datetime:    time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					},
					Station: Station{"BBC Radio Scotland", 9999},
				},
				{
					Match: Match{
						HomeTeam:    "Chelsea",
						AwayTeam:    "Spurs",
						Competition: "Premier League",
						Datetime:    time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					},
					Station: Radio5,
				},
			},
			output: []MatchDay{
				{
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								HomeTeam:    "Chelsea",
								AwayTeam:    "Tottenham",
								Competition: "Premier League",
								Datetime:    time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []Station{Talksport},
						},
					},
				},
				{
					DateTime: time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								HomeTeam:    "Chelsea",
								AwayTeam:    "Spurs",
								Competition: "Premier League",
								Datetime:    time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
							},
							Stations: []Station{Radio5},
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
		input  []MatchDay
		output CalData
	}{
		{
			input: []MatchDay{
				{
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []Listing{
						{
							Match: Match{
								HomeTeam:    "Chelsea",
								AwayTeam:    "Tottenham",
								Competition: "Premier League",
								Datetime:    time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []Station{Talksport, Radio5},
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
						Location: []Station{Talksport, Radio5},
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
