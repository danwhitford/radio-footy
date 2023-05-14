package feeds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/interchange"
)

func TestMapTeamNames(t *testing.T) {
	table := []struct {
		input  interchange.MergedMatch
		output interchange.MergedMatch
	}{
		{
			input: interchange.MergedMatch{
				Title: "Milan v Chelsea",
			},
			output: interchange.MergedMatch{
				Title: "AC Milan v Chelsea",
			},
		},
		{
			input: interchange.MergedMatch{
				Title: "Chelsea v Milan",
			},
			output: interchange.MergedMatch{
				Title: "Chelsea v AC Milan",
			},
		},
	}

	for _, tst := range table {
		mapTeamNames(&tst.input)
		if diff := cmp.Diff(tst.input, tst.output); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
