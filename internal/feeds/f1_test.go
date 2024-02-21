package feeds

import (
	_ "embed"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/urlgetter"
)

//go:embed f1_test_day.html
var f1Html string

func TestGetF1OnSky(t *testing.T) {
	want := []Broadcast{
		{
			Match: Match{
				Datetime:    time.Date(2024, 3, 2, 15, 0, 0, 0, time.UTC),
				HomeTeam:    "Gulf Air Bahrain Grand Prix",
				AwayTeam:    "",
				Competition: "F1",
			},
			Station: SkySports,
		},
	}

	if html == "" {
		t.Errorf("failed to embed file")
	}

	getter := urlgetter.StringGetter{
		Contents: f1Html,
	}
	got, err := getF1OnSky(getter)
	if err != nil {
		t.Fatalf("got error: %s", err)
	}

	if diff := cmp.Diff(want, got[:1]); diff != "" {
		t.Errorf("F1 mismatch (-want +got):\n%s", diff)
	}
}
