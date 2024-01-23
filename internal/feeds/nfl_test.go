package feeds

import (
	"testing"

	_ "embed"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/urlgetter"
)

//go:embed nfl_test_day.html
var html string

func TestGetNflOnSky(t *testing.T) {
	want := []Match{
		{
			Time:        "20:30",
			Date:        "Saturday, Jan 20",
			Stations:    []string{"Sky Sports"},
			Datetime:    "2024-01-20T20:30:00Z",
			HomeTeam:    "Baltimore Ravens",
			AwayTeam:    "Houston Texans",
			Competition: "NFL",
		},
		{
			Time:        "01:00",
			Date:        "Sunday, Jan 21",
			Stations:    []string{"Sky Sports"},
			Datetime:    "2024-01-21T01:00:00Z",
			HomeTeam:    "San Francisco 49ers",
			AwayTeam:    "Green Bay Packers",
			Competition: "NFL",
		},
		{
			Time:        "19:00",
			Date:        "Sunday, Jan 21",
			Stations:    []string{"Sky Sports"},
			Datetime:    "2024-01-21T19:00:00Z",
			HomeTeam:    "Detroit Lions",
			AwayTeam:    "Tampa Bay Buccaneers",
			Competition: "NFL",
		},
	}

	if html == "" {
		t.Errorf("failed to embed file")
	}

	getter := urlgetter.StringGetter{
		Contents: html,
	}
	got, err := getNflOnSky(getter)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NFL mismatch (-want +got):\n%s", diff)
	}
}
