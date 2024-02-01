package feeds

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"whitford.io/radiofooty/internal/urlgetter"
)

const englishFootballUrl = "https://www.live-footballontv.com/live-english-football-on-tv.html"

type tvFixture struct {
	homeTeam string
	awayTeam string
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
	"ITV4",
	"Channel 4",
	"TNT Sports 1",
	"TNT Sports 2",
	"TNT Sports 3",
	"TNT Sports Ultimate",
}

func getTvMatches(getter urlgetter.UrlGetter) ([]Broadcast, error) {
	re := regexp.MustCompile(`(\d+)(st|nd|rd|th)`)
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		return nil, fmt.Errorf("error loading location: %v", err)
	}

	html, err := getter.GetUrl(englishFootballUrl)
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
				teams = strings.TrimSpace(teams)
				splitTeams := strings.Split(teams, " v ")
				compName := div.Find("div", "class", "fixture__competition").Text()
				channels := div.Find("div", "class", "fixture__channel")
				channelStrings := make([]string, 0)
				for _, channelPill := range channels.FindAll("span", "class", "channel-pill") {
					channelString := channelPill.Text()

					if stringInSlice(channelString, channelsICareAbout) {
						if strings.HasPrefix(channelString, "Sky Sports") {
							channelString = "Sky Sports"
						}
						if strings.HasPrefix(channelString, "TNT Sports") {
							channelString = "TNT Sports"
						}
						if stringInSlice(channelString, channelStrings) {
							continue
						}
						channelStrings = append(channelStrings, channelString)
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

				fixtures = append(fixtures, tvFixture{
					homeTeam: splitTeams[0],
					awayTeam: splitTeams[1],
					compName: compName,
					dateTime: fixtureDateTime,
					channels: channelStrings,
				})
			}
		}
	}

	Matches := make([]Broadcast, 0)
	for _, fixture := range fixtures {
		for _, channel := range fixture.channels {
			Matches = append(Matches, Broadcast{
				Match: Match{
					HomeTeam:    fixture.homeTeam,
					AwayTeam:    fixture.awayTeam,
					Competition: fixture.compName,
					Datetime:    fixture.dateTime,
				},
				Station: channel,
			})
		}
	}
	return Matches, nil
}

func stringInSlice(str string, slice []string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
