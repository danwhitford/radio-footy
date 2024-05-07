package channel

import (
	_ "embed"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)
import "whitford.io/radiofooty/internal/broadcast"

//go:embed nfl_test_day.html
var nflHtml string

//go:embed f1_test_day.html
var f1Html string

type skyTestGetter struct{}

func (stg skyTestGetter) GetUrl(url string) ([]byte, error) {
	switch url {
	case "https://www.skysports.com/watch/nfl-on-sky":
		return []byte(nflHtml), nil
	case "https://www.skysports.com/watch/f1-on-sky":
		return []byte(f1Html), nil
	default:
		return []byte{}, fmt.Errorf("what is this url '%s'", url)
	}
}

func TestSkyGetMatches(t *testing.T) {
	want := []broadcast.Broadcast{
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 1, 20, 20, 30, 0, 0, time.UTC),
				HomeTeam:    "Baltimore Ravens",
				AwayTeam:    "Houston Texans",
				Competition: "NFL",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 1, 21, 0o1, 0o0, 0, 0, time.UTC),
				HomeTeam:    "San Francisco 49ers",
				AwayTeam:    "Green Bay Packers",
				Competition: "NFL",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 1, 21, 19, 0o0, 0, 0, time.UTC),
				HomeTeam:    "Detroit Lions",
				AwayTeam:    "Tampa Bay Buccaneers",
				Competition: "NFL",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 2, 29, 11, 30, 0, 0, time.UTC),
				HomeTeam:    "Gulf Air Bahrain Grand Prix - Practice",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 2, 29, 15, 0, 0, 0, time.UTC),
				HomeTeam:    "Gulf Air Bahrain Grand Prix - Practice",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 3, 1, 12, 30, 0, 0, time.UTC),
				HomeTeam:    "Gulf Air Bahrain Grand Prix - Practice",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 3, 1, 16, 0, 0, 0, time.UTC),
				HomeTeam:    "Gulf Air Bahrain Grand Prix - Qualifying",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 3, 1, 16, 22, 0, 0, time.UTC),
				HomeTeam:    "Gulf Air Bahrain Grand Prix - Qualifying",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 3, 1, 16, 45, 0, 0, time.UTC),
				HomeTeam:    "Gulf Air Bahrain Grand Prix - Qualifying",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 3, 2, 15, 0, 0, 0, time.UTC),
				HomeTeam:    "Gulf Air Bahrain Grand Prix",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 3, 7, 13, 30, 0, 0, time.UTC),
				HomeTeam:    "stc Saudi Arabian Grand Prix - Practice",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 3, 7, 17, 0, 0, 0, time.UTC),
				HomeTeam:    "stc Saudi Arabian Grand Prix - Practice",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: broadcast.SkySports,
		},
	}

	getter := skyTestGetter{}
	skyGetter := SkyGetter{getter}
	got, err := skyGetter.GetMatches()
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("Sky mismatch (-want +got):\n%s", diff)
	}

}
