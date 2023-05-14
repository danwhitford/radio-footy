package feeds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
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