package broadcast

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
		{
			"Bosnia & Herzegovina", "Bosnia-Herzegovina",
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
		{
			input:  "EFL Trophy Semi-Final",
			output: "EFL Trophy",
		},
		{
			input:  "EFL Trophy Final",
			output: "EFL Trophy",
		},
		{
			input:  "EFL Cup Football 2023-24",
			output: "EFL Cup",
		},
		{
			input:  "UEFA Europa Conference League Round of 16 1st Leg",
			output: "Europa Conference League",
		},
		{
			input:  "UEFA Europa League Round of 16 1st Leg",
			output: "Europa League",
		},
		{
			input:  "Europa League Football 2023-24",
			output: "Europa League",
		},
		{
			input:  "International Football 2023-24",
			output: "International Friendly",
		},
		{
			input:  "International Football 2024-25",
			output: "International Friendly",
		}, {
			input:  "International Friendly",
			output: "International Friendly",
		},
	}

	for _, tst := range table {
		got := mapCompName(tst.input)
		if diff := cmp.Diff(tst.output, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestTitle(t *testing.T) {
	table := []struct {
		match Match
		want  string
	}{
		{
			match: Match{
				Competition: "Premier League",
				HomeTeam:    "Bolton",
				AwayTeam:    "Barnsley",
			},
			want: "Bolton v Barnsley",
		},
		{
			match: Match{
				Competition: "NFL",
				HomeTeam:    "Cowboys",
				AwayTeam:    "Bills",
			},
			want: "Bills @ Cowboys",
		},
		{
			match: Match{
				Competition: "F1",
				HomeTeam:    "Silverstone Grand Prix",
				AwayTeam:    "",
			},
			want: "Silverstone Grand Prix",
		},
	}

	for _, tst := range table {
		got := tst.match.Title()
		if got != tst.want {
			t.Fatalf("wanted '%s' but got '%s'", tst.want, got)
		}
	}
}
