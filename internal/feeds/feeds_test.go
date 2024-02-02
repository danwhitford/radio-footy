package feeds

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMapTeamNames(t *testing.T) {
	table := []struct {
		input  Match
		output Match
	}{
		{
			input: Match{
				HomeTeam: "Milan",
				AwayTeam: "Chelsea",
			},
			output: Match{
				HomeTeam: "AC Milan",
				AwayTeam: "Chelsea",
			},
		},
		{
			input: Match{
				HomeTeam: "Chelsea",
				AwayTeam: "Milan",
			},
			output: Match{
				HomeTeam: "Chelsea",
				AwayTeam: "AC Milan",
			},
		},
	}

	for _, tst := range table {
		mapTeamNames(&tst.input)
		if diff := cmp.Diff(tst.output, tst.input); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestMapCompName(t *testing.T) {
	table := []struct {
		input  Match
		output Match
	}{
		{
			input: Match{
				Competition: "Premier League Football 2022-23",
			},
			output: Match{
				Competition: "Premier League",
			},
		},
	}

	for _, tst := range table {
		mapCompName(&tst.input)
		if diff := cmp.Diff(tst.output, tst.input); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestRollUpStations(t *testing.T) {
	table := []struct {
		input  []Broadcast
		output []Listing
	}{
		{
			input: []Broadcast{
				{
					Station: "BBC Radio 5",
					Match:   Match{},
				},
				{
					Station: "talkSPORT",
					Match:   Match{},
				},
			},
			output: []Listing{
				{
					Match:    Match{},
					Stations: []string{"talkSPORT", "BBC Radio 5"},
				},
			},
		},
		{
			input: []Broadcast{
				{
					Station: "BBC Radio 5",
					Match:   Match{},
				},
				{
					Station: "BBC Radio 5 Extra",
					Match:   Match{},
				},
			},
			output: []Listing{
				{
					Match:    Match{},
					Stations: []string{"BBC Radio 5", "BBC Radio 5 Extra"},
				},
			},
		},
		{
			input: []Broadcast{
				{
					Station: "talkSPORT",
					Match:   Match{},
				},
				{
					Station: "BBC Radio 5",
					Match:   Match{},
				},
				{
					Station: "Channel 4",
					Match:   Match{},
				},
			},
			output: []Listing{
				{
					Match:    Match{},
					Stations: []string{"Channel 4", "talkSPORT", "BBC Radio 5"},
				},
			},
		},
		{
			input: []Broadcast{
				{
					Station: "Radio 5 Sports Extra",
					Match: Match{
						HomeTeam: "West Ham United",
						AwayTeam: "AFC Bournemouth",
						Datetime: time.Date(2024, 2, 1, 19, 25, 0, 0, time.UTC),
					},
				},
				{
					Station: "TNT Sports",
					Match: Match{
						HomeTeam: "West Ham United",
						AwayTeam: "AFC Bournemouth",
						Datetime: time.Date(2024, 2, 1, 19, 30, 0, 0, time.UTC),
					},
				},
			},
			output: []Listing{
				{
					Stations: []string{"Radio 5 Sports Extra", "TNT Sports"},
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
		got := rollUpStations(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestRollUpDates(t *testing.T) {
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
		got := rollUpDates(tst.input)
		if diff := cmp.Diff(tst.output, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

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
		got := filterMatches(tst.input)
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
							Stations: []string{"talkSPORT"},
						},
						{
							Match: Match{

								HomeTeam: "Fulham",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []string{"talkSPORT2"},
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
							Stations: []string{"talkSPORT"},
						},
						{
							Match: Match{

								HomeTeam: "Fulham",
								AwayTeam: "Barnsley",
								Datetime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
							},
							Stations: []string{"talkSPORT2"},
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
		got := sortMatchDays(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
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
					Station: "talkSPORT",
				},
				{
					Match: Match{

						HomeTeam:    "Inverness",
						AwayTeam:    "Hibernian",
						Competition: "Scottish Premiership",
						Datetime:    time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					},
					Station: "BBC Radio Scotland",
				},
				{
					Match: Match{
						HomeTeam:    "Chelsea",
						AwayTeam:    "Spurs",
						Competition: "Premier League",
						Datetime:    time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					},
					Station: "BBC Radio 5",
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
							Stations: []string{"talkSPORT"},
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
							Stations: []string{"BBC Radio 5"},
						},
					},
				},
			},
		},
	}

	for _, tst := range table {
		got := MatchesToMatchDays(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestMatchDayToEventList(t *testing.T) {
	table := []struct {
		input  []MatchDay
		output []CalEvent
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
							Stations: []string{"talkSPORT", "BBC Radio 5"},
						},
					},
				},
			},
			output: []CalEvent{
				{
					Uid:      "chelseavtottenham/premierleague",
					DtStart:  "20210814T150000Z",
					Summary:  "Chelsea v Tottenham [Premier League]",
					Location: []string{"talkSPORT", "BBC Radio 5"},
				},
			},
		},
	}

	for _, tst := range table {
		got := MatchDayToEventList(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
