package feeds

import (
	"fmt"
	"regexp"
	"time"

	"github.com/anaskhan96/soup"
	"whitford.io/radiofooty/internal/filecacher"
)

const englishFootballUrl = "https://www.live-footballontv.com/live-english-football-on-tv.html"

type tvFixture struct {
	teams    string
	compName string
	dateTime time.Time
	channels []string
}

var channelsICareAbout = []string{
	"Sky Sports Football",
	"Sky Sports Main Event",
	"Sky Sports Premier League",
	"BBC One",
	"BBC Two",
	"Amazon Prime Video",
	"ITV1",
}

func getTvMatches() ([]MergedMatch, error) {
	re := regexp.MustCompile(`(\d+)(st|nd|rd|th)`)
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		return nil, fmt.Errorf("error loading location: %v", err)
	}

	getter := filecacher.NewHttpGetter()
	html, err := filecacher.GetUrl(englishFootballUrl, getter)
	if err != nil {
		return nil, err
	}
	root := soup.HTMLParse(string(html))

	fixtureGroups := root.FindAll("div", "class", "fixture-group")

	fixtures := make([]tvFixture, 0)

	for _, group := range fixtureGroups {
		var fixtureDate time.Time
		for _, div := range group.FindAll("div") {
			switch div.Attrs()["class"] {
			case "fixture-date":
				dateString := div.Text()
				// Clean string
				modifiedDateStr := re.ReplaceAllString(dateString, "$1")
				dateTime, err := time.Parse(
					"Monday 2 January 2006",
					modifiedDateStr,
				)
				if err != nil {
					return nil, err
				}
				fixtureDate = dateTime
			case "fixture":
				teams := div.Find("div", "class", "fixture__teams").Text()
				compName := div.Find("div", "class", "fixture__competition").Text()
				channels := div.Find("div", "class", "fixture__channel")
				channelStrings := make([]string, 0)
				for _, channelPill := range channels.FindAll("span", "class", "channel-pill") {
					channelString := channelPill.Text()
					if stringInSlice(channelString, channelsICareAbout) {
						channelStrings = append(channelStrings, channelPill.Text())
					}
				}
				if len(channelStrings) == 0 {
					continue
				}

				fixtureTime := div.Find("div", "class", "fixture__time").Text()
				var fixtureHours, fixtureMins int
				fmt.Sscanf(fixtureTime, "%d:%d", &fixtureHours, &fixtureMins)

				fixtureDateTime := time.Date(
					fixtureDate.Year(),
					fixtureDate.Month(),
					fixtureDate.Day(),
					fixtureHours,
					fixtureMins,
					0,
					0,
					loc,
				)
				// fixtureDateTime = fixtureDateTime.In(loc)

				fixtures = append(fixtures, tvFixture{
					teams:    teams,
					compName: compName,
					dateTime: fixtureDateTime,
					channels: channelStrings,
				})
			}
		}
	}

	mergedMatches := make([]MergedMatch, 0)
	for _, fixture := range fixtures {
		mergedMatches = append(mergedMatches, MergedMatch{
			Title:       fixture.teams,
			Competition: fixture.compName,
			Datetime:    fixture.dateTime.Format(time.RFC3339),
			Date:        fixture.dateTime.Format(niceDate),
			Time:        fixture.dateTime.Format(timeLayout),
			Stations:    fixture.channels,
		})
	}
	return mergedMatches, nil
}

func stringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
