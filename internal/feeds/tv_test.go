package feeds

import (
	"testing"
	"time"

	_ "embed"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/urlgetter"
)

//go:embed tv_test_day.html
var tvHtml string

func TestTv(t *testing.T) {
	want := Broadcast{
		Match: Match{
			Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
			HomeTeam:    "West Brom",
			AwayTeam:    "Chelsea",
			Competition: "FA Cup",
		},
		Station: ITV1,
	}

	getter := urlgetter.StringGetter{
		Contents: tvHtml,
	}
	got, err := getTvMatches(getter)
	if err != nil {
		t.Fatal(err)
	}
	for _, b := range got {
		if diff := cmp.Diff(want, b); diff == "" {
			return
		}
	}
	for _, b := range got {
		t.Logf("%#v\n", b)
	}
	t.Fatal("TV test game not found")
}
