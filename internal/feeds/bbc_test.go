package feeds

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBBCDayToMatch(t *testing.T) {
	table := []struct {
		input  BBCFeed
		output []Match
	}{
		{
			input: BBCFeed{
				Data: []BBCFeedData{
					{
						Data: []BBCProgramData{
							{
								Title: BBCTitles{
									Primary:   "5 Live Sport",
									Secondary: "Premier League Football 2022-23",
									Tertiary:  "Arsenal v Chelsea",
								},
								Start: "2020-12-26T17:30:00Z",
							},
						},
					},
				},
			},
			output: []Match{
				{
					HomeTeam:    "Arsenal",
					AwayTeam:    "Chelsea",
					Stations:    []string{"BBC Radio 5"},
					Competition: "Premier League Football 2022-23",
					Time:        "17:30",
					Date:        "Saturday, Dec 26",
					Datetime:    "2020-12-26T17:30:00Z",
				},
			},
		},
		{
			input: BBCFeed{
				Data: []BBCFeedData{
					{
						Data: []BBCProgramData{
							{
								Title: BBCTitles{
									Primary:   "5 Live Sport",
									Secondary: "International Football 2022-23",
									Tertiary:  "England Women v Portugal Women",
								},
								Start: "2023-07-01T14:15:00Z",
								Network: BBCNetwork{
									ShortTitle: "Radio 5 Live",
								},
								Synopses: BBCSynopses{
									Short: "Live football commentary of England Women v Portugal Women in an international friendly.",
								},
							},
						},
					},
				},
			},
			output: []Match{},
		},
		{
			input: BBCFeed{
				Data: []BBCFeedData{
					{
						Data: []BBCProgramData{
							{
								Title:    BBCTitles{Primary: "Cricket", Secondary: "Surrey v Nottinghamshire", Tertiary: ""},
								Start:    "2023-07-13T09:55:00Z",
								Network:  BBCNetwork{ShortTitle: "Radio 5 Sports Extra"},
								Synopses: BBCSynopses{Short: "Kevin Howells presents commentary of Surrey v Nottinghamshire in the County Championship."},
							},
						},
					},
				},
			},
			output: []Match{},
		},
		{
			input: BBCFeed{
				Data: []BBCFeedData{
					{
						Data: []BBCProgramData{
							{
								Title: BBCTitles{
									Primary:   "Cricket",
									Secondary: "Oval Invincibles Men v Welsh Fire Men",
									Tertiary:  "",
								},
								Start: "2023-08-06T18:00:00Z",
								Network: BBCNetwork{
									ShortTitle: "Radio 5 Sports Extra",
								},
								Synopses: BBCSynopses{
									Short: "Live commentary of Oval Invincibles Men v Welsh Fire Men in The Hundred at the Oval.",
								},
							},
						},
					},
				},
			},
			output: []Match{},
		},
	}

	for _, test := range table {
		got := bbcDayToMatches(test.input)
		if diff := cmp.Diff(test.output, got); diff != "" {
			t.Errorf("bbcDayToMatch() mismatch (-want +got):\n%s", diff)
		}
	}
}
