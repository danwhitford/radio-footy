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
		{
			input: Match{
				Competition: "International Friendlies",
			},
			output: Match{
				Competition: "International",
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
		input  []Match
		output []Match
	}{
		{
			input: []Match{
				{
					Stations: []string{"BBC Radio 5"},
					Time:     "17:30",
					Date:     "Saturday, Dec 26",
				},
				{
					Stations: []string{"talkSPORT"},
					Time:     "17:30",
					Date:     "Saturday, Dec 26",
				},
			},
			output: []Match{
				{
					Time:     "17:30",
					Date:     "Saturday, Dec 26",
					Stations: []string{"talkSPORT", "BBC Radio 5"},
				},
			},
		},
		{
			input: []Match{
				{
					Stations: []string{"BBC Radio 5"},
					Time:     "10:00",
					Date:     "Saturday, Dec 26",
				},
				{
					Stations: []string{"BBC Radio 5 Extra"},
					Time:     "12:00",
					Date:     "Saturday, Dec 26",
				},
			},
			output: []Match{
				{
					Time:     "10:00",
					Date:     "Saturday, Dec 26",
					Stations: []string{"BBC Radio 5", "BBC Radio 5 Extra"},
				},
			},
		},
		{
			input: []Match{
				{
					Stations: []string{"Sky Sports", "Sky Sports Football"},
					Time:     "10:00",
					Date:     "Saturday, Dec 26",
				},
				{
					Stations: []string{"BBC Radio 5"},
					Time:     "10:00",
					Date:     "Saturday, Dec 26",
				},
			},
			output: []Match{
				{
					Time:     "10:00",
					Date:     "Saturday, Dec 26",
					Stations: []string{"Sky Sports", "BBC Radio 5", "Sky Sports Football"},
				},
			},
		},
		{
			input: []Match{
				{
					Stations: []string{"talkSPORT"},
					Time:     "10:00",
					Date:     "Saturday, Dec 26",
				},
				{
					Stations: []string{"BBC Radio 5"},
					Time:     "10:00",
					Date:     "Saturday, Dec 26",
				},
				{
					Stations: []string{"Channel 4"},
					Time:     "10:00",
					Date:     "Saturday, Dec 26",
				},
			},
			output: []Match{
				{
					Time:     "10:00",
					Date:     "Saturday, Dec 26",
					Stations: []string{"talkSPORT", "BBC Radio 5", "Channel 4"},
				},
			},
		},
	}

	for _, tst := range table {
		got := rollUpStations(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestRollUpDates(t *testing.T) {
	table := []struct {
		input  []Match
		output []MatchDay
	}{
		{
			input: []Match{
				{
					Datetime: "2021-08-14T19:00:00Z",
					HomeTeam: "Milan",
					AwayTeam: "Chelsea",
				},
				{
					Datetime: "2021-08-14T21:00:00Z",
					HomeTeam: "Bolton",
					AwayTeam: "Barnsley",
				},
				{
					Datetime: "2021-08-15T19:00:00Z",
					HomeTeam: "Romsey",
					AwayTeam: "Worthing",
				},
			},
			output: []MatchDay{
				{
					NiceDate: "Saturday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Datetime: "2021-08-14T19:00:00Z",
							HomeTeam: "Milan",
							AwayTeam: "Chelsea",
						},
						{
							Datetime: "2021-08-14T21:00:00Z",
							HomeTeam: "Bolton",
							AwayTeam: "Barnsley",
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Datetime: "2021-08-15T19:00:00Z",
							HomeTeam: "Romsey",
							AwayTeam: "Worthing",
						},
					},
				},
			},
		},
	}

	for _, tst := range table {
		less := func(i, j MatchDay) bool {
			return i.NiceDate < j.NiceDate
		}
		got := rollUpDates(tst.input)
		if diff := cmp.Diff(tst.output, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestFilterMatches(t *testing.T) {
	table := []struct {
		input  []Match
		output []Match
	}{
		{
			input: []Match{
				{
					HomeTeam:    "Chelsea",
					AwayTeam:    "Milan",
					Competition: "Premier League",
				},
				{
					HomeTeam:    "Inverness",
					AwayTeam:    "Hibernian",
					Competition: "Scottish Premiership",
				},
			},
			output: []Match{
				{
					HomeTeam:    "Chelsea",
					AwayTeam:    "Milan",
					Competition: "Premier League",
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
					NiceDate: "Monday, Aug 18",
					DateTime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Date:     "Monday, Aug 18",
							HomeTeam: "Brentford",
							AwayTeam: "Arsenal",
							Datetime: "2021-08-18T19:00:00Z",
							Time:     "19:00",
						},
					},
				},
				{
					NiceDate: "Tuesday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Date:     "Tuesday, Aug 14",
							HomeTeam: "Bolton",
							AwayTeam: "Barnsley",
							Datetime: "2021-08-14T15:00:00Z",
							Time:     "15:00",
							Stations: []string{"talkSPORT"},
						},
						{
							Date:     "Tuesday, Aug 14",
							HomeTeam: "Fulham",
							AwayTeam: "Barnsley",
							Datetime: "2021-08-14T15:00:00Z",
							Time:     "15:00",
							Stations: []string{"talkSPORT2"},
						},
						{
							Date:     "Tuesday, Aug 14",
							HomeTeam: "Chelsea",
							AwayTeam: "Milan",
							Datetime: "2021-08-14T12:00:00Z",
							Time:     "12:00",
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Date:     "Wednesday, Aug 15",
							HomeTeam: "Romsey",
							AwayTeam: "Worthing",
							Datetime: "2021-08-15T15:00:00Z",
							Time:     "15:00",
						},
					},
				},
			},
			output: []MatchDay{
				{
					NiceDate: "Tuesday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Date:     "Tuesday, Aug 14",
							HomeTeam: "Chelsea",
							AwayTeam: "Milan",
							Datetime: "2021-08-14T12:00:00Z",
							Time:     "12:00",
						},
						{
							Date:     "Tuesday, Aug 14",
							HomeTeam: "Bolton",
							AwayTeam: "Barnsley",
							Datetime: "2021-08-14T15:00:00Z",
							Time:     "15:00",
							Stations: []string{"talkSPORT"},
						},
						{
							Date:     "Tuesday, Aug 14",
							HomeTeam: "Fulham",
							AwayTeam: "Barnsley",
							Datetime: "2021-08-14T15:00:00Z",
							Time:     "15:00",
							Stations: []string{"talkSPORT2"},
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Date:     "Wednesday, Aug 15",
							HomeTeam: "Romsey",
							AwayTeam: "Worthing",
							Datetime: "2021-08-15T15:00:00Z",
							Time:     "15:00",
						},
					},
				},
				{
					NiceDate: "Monday, Aug 18",
					DateTime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Date:     "Monday, Aug 18",
							HomeTeam: "Brentford",
							AwayTeam: "Arsenal",
							Datetime: "2021-08-18T19:00:00Z",
							Time:     "19:00",
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
		input  []Match
		output []MatchDay
	}{
		{
			input: []Match{
				{
					Date:        "Saturday, Aug 14",
					HomeTeam:    "Chelsea",
					AwayTeam:    "Tottenham",
					Competition: "Premier League",
					Datetime:    "2021-08-14T15:00:00Z",
					Stations:    []string{"talkSPORT", "BBC Radio 5"},
				},
				{
					Date:        "Saturday, Aug 14",
					HomeTeam:    "Inverness",
					AwayTeam:    "Hibernian",
					Competition: "Scottish Premiership",
					Datetime:    "2021-08-14T15:00:00Z",
					Stations:    []string{"BBC Radio Scotland"},
				},
				{
					Date:        "Sunday, Aug 15",
					HomeTeam:    "Chelsea",
					AwayTeam:    "Spurs",
					Competition: "Premier League",
					Datetime:    "2021-08-15T15:00:00Z",
					Stations:    []string{"BBC Radio 5"},
				},
			},
			output: []MatchDay{
				{
					NiceDate: "Saturday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Date:        "Saturday, Aug 14",
							HomeTeam:    "Chelsea",
							AwayTeam:    "Tottenham",
							Competition: "Premier League",
							Datetime:    "2021-08-14T15:00:00Z",
							Stations:    []string{"talkSPORT", "BBC Radio 5"},
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Date:        "Sunday, Aug 15",
							HomeTeam:    "Chelsea",
							AwayTeam:    "Spurs",
							Competition: "Premier League",
							Datetime:    "2021-08-15T15:00:00Z",
							Stations:    []string{"BBC Radio 5"},
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
					NiceDate: "Saturday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []Match{
						{
							Date:        "Saturday, Aug 14",
							HomeTeam:    "Chelsea",
							AwayTeam:    "Tottenham",
							Competition: "Premier League",
							Datetime:    "2021-08-14T15:00:00Z",
							Stations:    []string{"talkSPORT", "BBC Radio 5"},
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
