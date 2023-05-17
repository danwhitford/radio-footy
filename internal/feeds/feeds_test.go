package feeds

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestMapTeamNames(t *testing.T) {
	table := []struct {
		input  MergedMatch
		output MergedMatch
	}{
		{
			input: MergedMatch{
				Title: "Milan v Chelsea",
			},
			output: MergedMatch{
				Title: "AC Milan v Chelsea",
			},
		},
		{
			input: MergedMatch{
				Title: "Chelsea v Milan",
			},
			output: MergedMatch{
				Title: "Chelsea v AC Milan",
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
		input  MergedMatch
		output MergedMatch
	}{
		{
			input: MergedMatch{
				Competition: "Premier League Football 2022-23",
			},
			output: MergedMatch{
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
		input  []MergedMatch
		output []MergedMatch
	}{
		{
			input: []MergedMatch{
				{
					Stations: []string{"BBC Radio 5"},
				},
				{
					Stations: []string{"talkSPORT"},
				},
			},
			output: []MergedMatch{
				{
					Stations: []string{"talkSPORT", "BBC Radio 5"},
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
		input  []MergedMatch
		output []MergedMatchDay
	}{
		{
			input: []MergedMatch{
				{
					Datetime: "2021-08-14T19:00:00Z",
					Title:    "Chelsea v Milan",
				},
				{
					Datetime: "2021-08-14T21:00:00Z",
					Title:    "Bolton v Barnslay",
				},
				{
					Datetime: "2021-08-15T19:00:00Z",
					Title:    "Romsey v Worthing",
				},
			},
			output: []MergedMatchDay{
				{
					NiceDate: "Saturday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Datetime: "2021-08-14T19:00:00Z",
							Title:    "Chelsea v Milan",
						},
						{
							Datetime: "2021-08-14T21:00:00Z",
							Title:    "Bolton v Barnslay",
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Datetime: "2021-08-15T19:00:00Z",
							Title:    "Romsey v Worthing",
						},
					},
				},
			},
		},
	}

	for _, tst := range table {
		less := func(i, j MergedMatchDay) bool {
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
		input  []MergedMatch
		output []MergedMatch
	}{
		{
			input: []MergedMatch{
				{
					Title:       "Chelsea v Milan",
					Competition: "Premier League",
				},
				{
					Title:       "Inverness v Hibernian",
					Competition: "Scottish Premiership",
				},
			},
			output: []MergedMatch{
				{
					Title:       "Chelsea v Milan",
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
		input  []MergedMatchDay
		output []MergedMatchDay
	}{
		{
			input: []MergedMatchDay{
				{
					NiceDate: "Monday, Aug 18",
					DateTime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Date:     "Monday, Aug 18",
							Title:    "Brentford v Arsenal",
							Datetime: "2021-08-18T19:00:00Z",
							Time:     "19:00",
						},
					},
				},
				{
					NiceDate: "Tuesday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Date:     "Tuesday, Aug 14",
							Title:    "Chelsea v Milan",
							Datetime: "2021-08-14T15:00:00Z",
							Time:     "15:00",
							Stations: []string{"talkSPORT"},
						},
						{
							Date:     "Tuesday, Aug 14",
							Title:    "Fulham v Brentford",
							Datetime: "2021-08-14T15:00:00Z",
							Time:     "15:00",
							Stations: []string{"talkSPORT2"},
						},
						{
							Date:     "Tuesday, Aug 14",
							Title:    "Bolton v Barnslay",
							Datetime: "2021-08-14T12:00:00Z",
							Time:     "12:00",
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Date:     "Wednesday, Aug 15",
							Title:    "Romsey v Worthing",
							Datetime: "2021-08-15T15:00:00Z",
							Time:     "15:00",
						},
					},
				},
			},
			output: []MergedMatchDay{
				{
					NiceDate: "Tuesday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Date:     "Tuesday, Aug 14",
							Title:    "Bolton v Barnslay",
							Datetime: "2021-08-14T12:00:00Z",
							Time:     "12:00",
						},
						{
							Date:     "Tuesday, Aug 14",
							Title:    "Chelsea v Milan",
							Datetime: "2021-08-14T15:00:00Z",
							Time:     "15:00",
							Stations: []string{"talkSPORT"},
						},
						{
							Date:     "Tuesday, Aug 14",
							Title:    "Fulham v Brentford",
							Datetime: "2021-08-14T15:00:00Z",
							Time:     "15:00",
							Stations: []string{"talkSPORT2"},
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Date:     "Wednesday, Aug 15",
							Title:    "Romsey v Worthing",
							Datetime: "2021-08-15T15:00:00Z",
							Time:     "15:00",
						},
					},
				},
				{
					NiceDate: "Monday, Aug 18",
					DateTime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Date:     "Monday, Aug 18",
							Title:    "Brentford v Arsenal",
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

func TestFuzzyMergeTeams(t *testing.T) {
	table := []struct {
		input  []MergedMatch
		output []MergedMatch
	}{
		{
			input: []MergedMatch{
				{
					Title:       "Chelsea v Tottenham",
					Competition: "Premier League",
					Datetime:    "2021-08-14T15:00:00Z",
					Stations:    []string{"talkSPORT"},
				},
				{
					Title:       "Chelsea v Spurs",
					Competition: "Premier League",
					Datetime:    "2021-08-14T15:00:00Z",
					Stations:    []string{"talkSPORT2"},
				},
				{
					Title:       "Inverness v Hibernian",
					Competition: "Scottish Premiership",
					Datetime:    "2021-08-14T15:00:00Z",
					Stations:    []string{"BBC Radio Scotland"},
				},
			},
			output: []MergedMatch{
				{
					Title:       "Chelsea v Tottenham",
					Competition: "Premier League",
					Datetime:    "2021-08-14T15:00:00Z",
					Stations:    []string{"talkSPORT", "talkSPORT2"},
				},
				{
					Title:       "Inverness v Hibernian",
					Competition: "Scottish Premiership",
					Datetime:    "2021-08-14T15:00:00Z",
					Stations:    []string{"BBC Radio Scotland"},
				},
			},
		},
	}

	for _, tst := range table {
		got := fuzzyMergeTeams(tst.input)
		less := func(i, j MergedMatch) bool {
			return i.Title < j.Title
		}
		if diff := cmp.Diff(tst.output, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestMergedMatchesToMergedMatchDays(t *testing.T) {
	table := []struct {
		input  []MergedMatch
		output []MergedMatchDay
	}{
		{
			input: []MergedMatch{
				{
					Date:        "Saturday, Aug 14",
					Title:       "Chelsea v Tottenham",
					Competition: "Premier League",
					Datetime:    "2021-08-14T15:00:00Z",
					Stations:    []string{"talkSPORT", "BBC Radio 5"},
				},
				{
					Date:        "Saturday, Aug 14",
					Title:       "Inverness v Hibernian",
					Competition: "Scottish Premiership",
					Datetime:    "2021-08-14T15:00:00Z",
					Stations:    []string{"BBC Radio Scotland"},
				},
				{
					Date:        "Sunday, Aug 15",
					Title:       "Chelsea v Spurs",
					Competition: "Premier League",
					Datetime:    "2021-08-15T15:00:00Z",
					Stations:    []string{"BBC Radio 5"},
				},
			},
			output: []MergedMatchDay{
				{
					NiceDate: "Saturday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Date:        "Saturday, Aug 14",
							Title:       "Chelsea v Tottenham",
							Competition: "Premier League",
							Datetime:    "2021-08-14T15:00:00Z",
							Stations:    []string{"talkSPORT", "BBC Radio 5"},
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Date:        "Sunday, Aug 15",
							Title:       "Chelsea v Spurs",
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
		got := mergedMatchesToMergedMatchDays(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestMergedMatchDayToEventList(t *testing.T) {
	table := []struct {
		input  []MergedMatchDay
		output []CalEvent
	}{
		{
			input: []MergedMatchDay{
				{
					NiceDate: "Saturday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []MergedMatch{
						{
							Date:        "Saturday, Aug 14",
							Title:       "Chelsea v Tottenham",
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
		got := MergedMatchDayToEventList(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
