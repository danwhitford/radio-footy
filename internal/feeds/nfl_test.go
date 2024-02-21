package feeds

import (
	_ "embed"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/urlgetter"
)

//go:embed nfl_test_day.html
var html string

func TestGetNflOnSky(t *testing.T) {
	want := []Broadcast{
		{
			Match: Match{
				Datetime:    time.Date(2024, 1, 20, 20, 30, 0, 0, time.UTC),
				HomeTeam:    "Baltimore Ravens",
				AwayTeam:    "Houston Texans",
				Competition: "NFL",
			},
			Station: SkySports,
		},
		{
			Match: Match{
				Datetime:    time.Date(2024, 1, 21, 0o1, 0o0, 0, 0, time.UTC),
				HomeTeam:    "San Francisco 49ers",
				AwayTeam:    "Green Bay Packers",
				Competition: "NFL",
			},
			Station: SkySports,
		},
		{
			Match: Match{
				Datetime:    time.Date(2024, 1, 21, 19, 0o0, 0, 0, time.UTC),
				HomeTeam:    "Detroit Lions",
				AwayTeam:    "Tampa Bay Buccaneers",
				Competition: "NFL",
			},
			Station: SkySports,
		},
	}

	if html == "" {
		t.Errorf("failed to embed file")
	}

	getter := urlgetter.StringGetter{
		Contents: html,
	}
	nflGetter := nflMatchGetter{getter}
	got, err := nflGetter.getMatches()
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NFL mismatch (-want +got):\n%s", diff)
	}
}
