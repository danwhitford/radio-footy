package channel

import (
	_ "embed"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/broadcast"
	"whitford.io/radiofooty/internal/urlgetter"
)

//go:embed tv_test_day.html
var tvHtml string

func TestTv(t *testing.T) {
	want := broadcast.Broadcast{
		Match: broadcast.Match{
			Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
			HomeTeam:    "West Brom",
			AwayTeam:    "Chelsea",
			Competition: "FA Cup",
		},
		Station: broadcast.ITV1,
	}

	getter := urlgetter.StringGetter{
		Contents: tvHtml,
	}
	tmg := TvMatchGetter{getter}
	got, err := tmg.GetMatches()
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
