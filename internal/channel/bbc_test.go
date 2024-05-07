package channel

import (
	"sort"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"whitford.io/radiofooty/internal/broadcast"
)

func TestBBCDayToMatch(t *testing.T) {
	table := []struct {
		input  BBCFeed
		output []broadcast.Broadcast
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
			output: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "Chelsea",
						Competition: "Premier League",
						Datetime:    time.Date(2020, 12, 26, 17, 30, 0, 0, time.UTC),
					},
					Station: broadcast.BlankStation,
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
			output: []broadcast.Broadcast{},
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
			output: []broadcast.Broadcast{},
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
			output: []broadcast.Broadcast{},
		},
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
								Network: BBCNetwork{
									ShortTitle: "Radio 5 Sports Extra",
								},
								Start: "2020-12-26T17:30:00Z",
							},
						},
					},
				},
			},
			output: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "Chelsea",
						Competition: "Premier League",
						Datetime:    time.Date(2020, 12, 26, 17, 30, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5Extra,
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
									Primary:   "Six Nations 2024",
									Secondary: "Italy v England",
								},
								Start: "2024-02-03T14:15:00Z",
								Network: BBCNetwork{
									ShortTitle: "Radio 5 Sports Extra",
								},
								Synopses: BBCSynopses{
									Short: "Live rugby union commentary of Italy v England in the Six Nations at Stadio Olimpico.",
								},
							},
						},
					},
					{
						Data: []BBCProgramData{
							{Title: BBCTitles{Primary: "Six Nations 2024", Secondary: "Round two preview", Tertiary: ""}, Start: "2024-02-08T19:00:00Z", Network: BBCNetwork{ShortTitle: "Radio 5 Live"}, Synopses: BBCSynopses{Short: "Sonja McLaughlan looks ahead to round two of the Six Nations Championships."}},
						},
					},
				},
			},
			output: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Italy",
						AwayTeam:    "England",
						Competition: "Six Nations",
						Datetime:    time.Date(2024, 2, 3, 14, 15, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5Extra,
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
									Secondary: "Premier League Football 2023-24",
									Tertiary:  "Arsenal v West Bromwich Albion",
								},
								Network: BBCNetwork{
									ShortTitle: "Radio 5 Sports Extra",
								},
								Start: "2024-02-07T20:00:00Z",
							},
						},
					},
				},
			},
			output: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5Extra,
				},
			},
		},
	}

	for _, test := range table {
		got, err := bbcDayToMatches(test.input)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(test.output, got); diff != "" {
			t.Fatalf("bbcDayToMatch() mismatch (-want +got):\n%s", diff)
		}
	}
}

func TestDeupeBbcMatches(t *testing.T) {
	table := []struct {
		in   []broadcast.Broadcast
		want []broadcast.Broadcast
	}{
		{
			in: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5Extra,
				},
			},
			want: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5Extra,
				},
			},
		},
		{
			in: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5,
				},
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5Extra,
				},
			},
			want: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5,
				},
			},
		},
		{
			in: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 18, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5,
				},
				{
					Match: broadcast.Match{
						HomeTeam:    "Liverpool",
						AwayTeam:    "Everton",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5Extra,
				},
			},
			want: []broadcast.Broadcast{
				{
					Match: broadcast.Match{
						HomeTeam:    "Arsenal",
						AwayTeam:    "West Brom",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 18, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5,
				},
				{
					Match: broadcast.Match{
						HomeTeam:    "Liverpool",
						AwayTeam:    "Everton",
						Competition: "Premier League",
						Datetime:    time.Date(2024, 2, 7, 20, 0, 0, 0, time.UTC),
					},
					Station: broadcast.Radio5Extra,
				},
			},
		},
	}

	trans := cmp.Transformer("Sort", func(in []broadcast.Broadcast) []broadcast.Broadcast {
		out := append([]broadcast.Broadcast(nil), in...) // Copy input to avoid mutating it
		sort.Slice(out, func(i, j int) bool {
			return out[i].Datetime.Before(out[j].Datetime)
		})
		return out
	})

	for _, tst := range table {
		got := dedupeBbcMatches(tst.in)

		if diff := cmp.Diff(tst.want, got, trans); diff != "" {
			t.Fatalf("dedupeBbcMatches() mismatch (-want +got):\n%s", diff)
		}
	}
}
