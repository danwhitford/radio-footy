package feeds

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

var rugbyText = `Friday 2 February
France v Ireland (8:00pm) ITV1
Saturday 3 February
Italy v England (2:15pm) ITV1
Wales v Scotland (4:45pm) BBC One
Saturday 10 February
Scotland v France (2:15pm) BBC One
England v Wales (4:45pm) ITV1
Sunday 11 February
Ireland v Italy (3:00pm) ITV1
Saturday 24 February
Ireland v Wales (2:15pm) ITV1
Scotland v England (4:45pm) BBC One
Sunday 25 February
France v Italy (3:00pm) ITV1
Saturday 9 March
Italy v Scotland (2:15pm) ITV1
England v Ireland (4:45pm) ITV1
Sunday 10 March
Wales v France (3:00pm) BBC One
Saturday 16 March
Wales v Italy (2:15pm) BBC One
Ireland v Scotland (4:45pm) ITV1
France v England (8:00pm) ITV1
`

func parseRugby() ([]Broadcast, error) {
	broadcasts := make([]Broadcast, 0)
	sr := strings.NewReader(rugbyText)

	var date time.Time
	scanner := bufio.NewScanner(sr)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "Friday") ||
			strings.HasPrefix(line, "Saturday") ||
			strings.HasPrefix(line, "Sunday") {
			d, err := time.Parse("Monday 2 January", line)
			if err != nil {
				return broadcasts, err
			}
			date = d
		} else {
			var hTeam, aTeam, channel string
			var hours, minutes int
			fmt.Sscanf(line, "%s v %s (%d:%dpm) %s", &hTeam, &aTeam, &hours, &minutes, &channel)
			if channel == "BBC" {
				channel = "BBC One"
			}

			broadcasts = append(broadcasts, Broadcast{
				Match: Match{
					HomeTeam:    hTeam,
					AwayTeam:    aTeam,
					Competition: "Six Nations 2024",
					Datetime: time.Date(
						2024,
						date.Month(),
						date.Day(),
						hours+12,
						minutes,
						0,
						0,
						time.UTC,
					),
				},
				Station: StationFromString(channel),
			})
		}
	}

	return broadcasts, nil
}

func cricket() []Broadcast {
	gameDates := []struct {
		from time.Time
		to   time.Time
	}{
		{
			time.Date(2024, 1, 25, 4, 0, 0, 0, time.UTC),
			time.Date(2024, 1, 29, 12, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2024, 2, 2, 4, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 6, 12, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2024, 2, 15, 4, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 19, 12, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2024, 2, 23, 4, 0, 0, 0, time.UTC),
			time.Date(2024, 2, 27, 12, 0, 0, 0, time.UTC),
		},
		{
			time.Date(2024, 3, 7, 4, 0, 0, 0, time.UTC),
			time.Date(2024, 3, 11, 12, 0, 0, 0, time.UTC),
		},
	}

	broadcasts := make([]Broadcast, 0)
	for _, date := range gameDates {
		day := date.from
		for day.Before(date.to) {
			broadcasts = append(broadcasts, Broadcast{
				Match: Match{
					HomeTeam:    "India",
					AwayTeam:    "England",
					Competition: "Test Match",
					Datetime:    day,
				},
				Station: Talksport2,
			})
			day = day.AddDate(0, 0, 1)
		}
	}
	return broadcasts
}

func filterOld(bb []Broadcast) []Broadcast {
	out := make([]Broadcast, 0)
	for _, b := range bb {
		if b.Datetime.After(time.Now()) {
			out = append(out, b)
		}
	}
	return out
}

type manualGetter struct{}

func (mg manualGetter) getMatches() ([]Broadcast, error) {
	rugby, err := parseRugby()
	if err != nil {
		return rugby, err
	}
	crick := cricket()

	games := make([]Broadcast, 0)
	games = append(games, rugby...)
	games = append(games, crick...)

	return filterOld(games), nil
}
