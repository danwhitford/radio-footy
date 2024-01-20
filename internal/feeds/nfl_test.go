package feeds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	_ "embed"
)

//go:embed nfl_test_day.html
var html string

func TestGetNflOnSky(t *testing.T) {
	want := []MergedMatch{
		{
			Time: "20:30",
			Date: "Saturday, Jan 20",
			Stations: []string{"Sky Sports"},
			Datetime: "2024-01-20T20:30:00Z",
			Title: "Houston Texans v Baltimore Ravens",
			Competition: "NFL",
		},
		{
			Time: "01:00",
			Date: "Sunday, Jan 21",
			Stations: []string{"Sky Sports"},
			Datetime: "2024-01-21T01:00:00Z",
			Title: "Green Bay Packers v San Francisco 49ers",
			Competition: "NFL",
		},
		{
			Time: "19:00",
			Date: "Sunday, Jan 21",
			Stations: []string{"Sky Sports"},
			Datetime: "2024-01-21T19:00:00Z",
			Title: "Tampa Bay Buccaneers v Detroit Lions",
			Competition: "NFL",
		},
	}

	if html == "" {
		t.Errorf("failed to embed file")
	}
	
	got, err := nflPageToMergedMatches(html)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("NFL mismatch (-want +got):\n%s", diff)
	}
}
