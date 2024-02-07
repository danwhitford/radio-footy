package feeds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMapTeamNames(t *testing.T) {
	table := []struct {
		input  string
		output string
	}{
		{
			"Milan", "AC Milan",
		},
		{
			"Chelsea", "Chelsea",
		},
	}

	for _, tst := range table {
		got := mapTeamName(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestMapCompName(t *testing.T) {
	table := []struct {
		input  string
		output string
	}{
		{
			input:  "Premier League Football 2022-23",
			output: "Premier League",
		},
		{
			input:  "UEFA Champions League",
			output: "Champions League",
		},
		{
			input:  "UEFA Champions League┬áRound of 16 1st Leg",
			output: "Champions League",
		},
		{
			input:  "Champions League",
			output: "Champions League",
		},
	}

	for _, tst := range table {
		got := mapCompName(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
