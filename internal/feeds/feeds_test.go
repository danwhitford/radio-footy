package feeds

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"whitford.io/radiofooty/internal/interchange"
)

func TestMapTeamNames(t *testing.T) {
	table := []struct {
		input  interchange.MergedMatch
		output interchange.MergedMatch
	}{
		{
			input: interchange.MergedMatch{
				Title: "Milan v Chelsea",
			},
			output: interchange.MergedMatch{
				Title: "AC Milan v Chelsea",
			},
		},
		{
			input: interchange.MergedMatch{
				Title: "Chelsea v Milan",
			},
			output: interchange.MergedMatch{
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
		input  interchange.MergedMatch
		output interchange.MergedMatch
	}{
		{
			input: interchange.MergedMatch{
				Competition: "Premier League Football 2022-23",
			},
			output: interchange.MergedMatch{
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
		input  []interchange.MergedMatch
		output []interchange.MergedMatch
	}{
		{
			input: []interchange.MergedMatch{
				{
					Station: "BBC Radio 5 Live",
				},
				{
					Station: "talkSPORT",
				},
			},
			output: []interchange.MergedMatch{
				{
					Station: "talkSPORT | BBC Radio 5 Live",
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
		input  []interchange.MergedMatch
		output []interchange.MergedMatchDay
	}{
		{
			input: []interchange.MergedMatch{
				{
					Date:     "Saturday, Aug 14",
					Datetime: "2021-08-14T19:00:00Z",
					Title:    "Chelsea v Milan",
				},
				{
					Date:     "Saturday, Aug 14",
					Datetime: "2021-08-14T21:00:00Z",
					Title:    "Bolton v Barnslay",
				},
				{
					Date:     "Sunday, Aug 15",
					Datetime: "2021-08-15T19:00:00Z",
					Title:    "Romsey v Worthing",
				},
			},
			output: []interchange.MergedMatchDay{
				{
					NiceDate: "Saturday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 0, 0, 0, 0, time.UTC),
					Matches: []interchange.MergedMatch{
						{
							Date:     "Saturday, Aug 14",
							Datetime: "2021-08-14T19:00:00Z",
							Title:    "Chelsea v Milan",
						},
						{
							Date:     "Saturday, Aug 14",
							Datetime: "2021-08-14T21:00:00Z",
							Title:    "Bolton v Barnslay",
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 0, 0, 0, 0, time.UTC),
					Matches: []interchange.MergedMatch{
						{
							Date:     "Sunday, Aug 15",
							Datetime: "2021-08-15T19:00:00Z",
							Title:    "Romsey v Worthing",
						},
					},
				},
			},
		},
	}

	for _, tst := range table {
		got := rollUpDates(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestFilterMatches(t *testing.T) {
	table := []struct {
		input  []interchange.MergedMatch
		output []interchange.MergedMatch
	}{
		{
			input: []interchange.MergedMatch{
				{
					Title:       "Chelsea v Milan",
					Competition: "Premier League",
				},
				{
					Title:       "Inverness v Hibernian",
					Competition: "Scottish Premiership",
				},
			},
			output: []interchange.MergedMatch{
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
		input  []interchange.MergedMatchDay
		output []interchange.MergedMatchDay
	}{
		{
			input: []interchange.MergedMatchDay{
				{
					NiceDate: "Monday, Aug 18",
					DateTime: time.Date(2021, 8, 18, 19, 0, 0, 0, time.UTC),
					Matches: []interchange.MergedMatch{
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
					Matches: []interchange.MergedMatch{
						{
							Date:     "Tuesday, Aug 14",
							Title:    "Chelsea v Milan",
							Datetime: "2021-08-14T15:00:00Z",
							Time:     "15:00",
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
					Matches: []interchange.MergedMatch{
						{
							Date:     "Wednesday, Aug 15",
							Title:    "Romsey v Worthing",
							Datetime: "2021-08-15T15:00:00Z",
							Time:     "15:00",
						},
					},
				},
			},
			output: []interchange.MergedMatchDay{
				{
					NiceDate: "Tuesday, Aug 14",
					DateTime: time.Date(2021, 8, 14, 15, 0, 0, 0, time.UTC),
					Matches: []interchange.MergedMatch{
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
						},
					},
				},
				{
					NiceDate: "Sunday, Aug 15",
					DateTime: time.Date(2021, 8, 15, 15, 0, 0, 0, time.UTC),
					Matches: []interchange.MergedMatch{
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
					Matches: []interchange.MergedMatch{
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
		input  []interchange.MergedMatch
		output []interchange.MergedMatch
	}{
		{
			input: []interchange.MergedMatch{
				{
					Title:       "Chelsea v Tottenham",
					Competition: "Premier League",
					Datetime:    "2021-08-14T15:00:00Z",
					Station:     "talkSPORT",
				},
				{
					Title:       "Chelsea v Spurs",
					Competition: "Premier League",
					Datetime:    "2021-08-14T15:00:00Z",
					Station:     "BBC Radio 5 Live",
				},
				{
					Title:       "Inverness v Hibernian",
					Competition: "Scottish Premiership",
					Datetime:    "2021-08-14T15:00:00Z",
					Station:     "BBC Radio Scotland",
				},
			},
			output: []interchange.MergedMatch{
				{
					Title:       "Chelsea v Tottenham",
					Competition: "Premier League",
					Datetime:    "2021-08-14T15:00:00Z",
					Station:     "talkSPORT | BBC Radio 5 Live",
				},
				{
					Title:       "Inverness v Hibernian",
					Competition: "Scottish Premiership",
					Datetime:    "2021-08-14T15:00:00Z",
					Station:     "BBC Radio Scotland",
				},
			},
		},
	}

	for _, tst := range table {
		got := fuzzyMergeTeams(tst.input)
		less := func(i, j int) bool {
			return got[i].Title < got[j].Title
		}
		if diff := cmp.Diff(tst.output, got, cmpopts.SortSlices(less)); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
