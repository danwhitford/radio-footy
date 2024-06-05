package channel

import (
	_ "embed"
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)
import "whitford.io/radiofooty/internal/broadcast"

//go:embed sky_tests/nfl_test_day.html
var nflHtml string

//go:embed sky_tests/f1_test_day.html
var f1Html string

//go:embed sky_tests/cricket_test_day.html
var cricketHtml string

type skyTestGetter struct{}

func (stg skyTestGetter) GetUrl(url string) ([]byte, error) {
	switch url {
	case "https://www.skysports.com/watch/nfl-on-sky":
		return []byte(nflHtml), nil
	case "https://www.skysports.com/watch/f1-on-sky":
		return []byte(f1Html), nil
	case "https://www.skysports.com/watch/cricket-on-sky":
		return []byte(cricketHtml), nil
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
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 5, 14, 0, 0, 0, time.UTC),
				HomeTeam:    "India",
				AwayTeam:    "Ireland",
				Competition: "ICC Men's T20 World Cup",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 5, 23, 0, 0, 0, time.UTC),
				HomeTeam:    "Papua New Guinea",
				AwayTeam:    "Uganda",
				Competition: "ICC Men's T20 World Cup",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 6, 0, 0, 0, 0, time.UTC),
				HomeTeam:    "Australia",
				AwayTeam:    "Oman",
				Competition: "ICC Men's T20 World Cup",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 6, 15, 0, 0, 0, time.UTC),
				HomeTeam:    "USA",
				AwayTeam:    "Pakistan",
				Competition: "ICC Men's T20 World Cup",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 6, 18, 30, 0, 0, time.UTC),
				HomeTeam:    "Namibia",
				AwayTeam:    "Scotland",
				Competition: "ICC Men's T20 World Cup",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 7, 14, 0, 0, 0, time.UTC),
				HomeTeam:    "Canada",
				AwayTeam:    "Ireland",
				Competition: "ICC Men's T20 World Cup",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 7, 18, 0, 0, 0, time.UTC),
				HomeTeam:    "Northampton",
				AwayTeam:    "Worcester",
				Competition: "T20",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 7, 23, 0, 0, 0, time.UTC),
				HomeTeam:    "New Zealand",
				AwayTeam:    "Afghanistan",
				Competition: "ICC Men's T20 World Cup",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 8, 0, 0, 0, 0, time.UTC),
				HomeTeam:    "Sri Lanka",
				AwayTeam:    "Bangladesh",
				Competition: "ICC Men's T20 World Cup",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 8, 14, 0, 0, 0, time.UTC),
				HomeTeam:    "Netherlands",
				AwayTeam:    "South Africa",
				Competition: "ICC Men's T20 World Cup",
			},
			Station: broadcast.SkySports,
		},
		{
			Match: broadcast.Match{
				Datetime:    time.Date(2024, 6, 8, 16, 30, 0, 0, time.UTC),
				HomeTeam:    "Australia",
				AwayTeam:    "England",
				Competition: "ICC Men's T20 World Cup",
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
