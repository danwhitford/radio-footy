package channel

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"whitford.io/radiofooty/internal/broadcast"
	"whitford.io/radiofooty/internal/urlgetter"
)

type SkyGetter struct {
	Urlgetter urlgetter.UrlGetter
}

type skyPage struct {
	competition   string
	teamExtractor teamExtractor
	url           string
}

type teamExtractor func([]soup.Root) (string, string, bool)

var dateRe *regexp.Regexp = regexp.MustCompile(`(\d+)(st|nd|rd|th)`)

func (sg SkyGetter) GetMatches() ([]broadcast.Broadcast, error) {
	pages := []skyPage{
		{
			"NFL",
			func(eventTitles []soup.Root) (string, string, bool) {
				homeTeam := eventTitles[1].Text()
				awayTeam := eventTitles[0].Text()
				return homeTeam, awayTeam, true
			},
			"https://www.skysports.com/watch/nfl-on-sky",
		},
		{
			"F1",
			func(eventTitles []soup.Root) (string, string, bool) {
				if strings.HasSuffix(eventTitles[0].Text(), " - Pit Lane Live") {
					return "", "", false
				}
				raceTitle := eventTitles[0].Text()
				raceTitle, _ = strings.CutSuffix(raceTitle, " - Race")
				raceTitle, _ = strings.CutSuffix(raceTitle, " 1")
				raceTitle, _ = strings.CutSuffix(raceTitle, " 2")
				raceTitle, _ = strings.CutSuffix(raceTitle, " 3")
				return raceTitle, "", true
			},
			"https://www.skysports.com/watch/f1-on-sky",
		},
	}

	broadcasts := make([]broadcast.Broadcast, 0)
	for _, page := range pages {
		html, err := sg.Urlgetter.GetUrl(page.url)
		if err != nil {
			return nil, err
		}
		bb, err := skyPageToMatches(string(html), page.competition, page.teamExtractor)
		if err != nil {
			return broadcasts, err
		}
		broadcasts = append(broadcasts, bb...)
	}

	return broadcasts, nil
}

func skyPageToMatches(html, comp string, extr teamExtractor) ([]broadcast.Broadcast, error) {
	Matches := make([]broadcast.Broadcast, 0)

	re := regexp.MustCompile(`\([0-9:]+\)`)
	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		return nil, fmt.Errorf("error loading location: %v", err)
	}
	root := soup.HTMLParse(html)

	var pageData soup.Root
	boxes := root.FindAll("div", "class", "box")
	for _, pd := range boxes {
		_, prs := pd.Attrs()["data-current-page"]
		if prs {
			pageData = pd
			break
		}
	}

	if pageData.Pointer == nil {
		return Matches, nil
	}
	var curDate time.Time
	for _, child := range pageData.Children() {
		switch child.NodeValue {
		case "h3":
			{
				datestring := child.Text()
				cleanstring := dateRe.ReplaceAllString(datestring, "$1")
				date, err := time.Parse("Mon 2 January", cleanstring)
				if err != nil {
					return Matches, fmt.Errorf("could not parse %s: %s", cleanstring, err)
				}
				curDate = time.Date(
					time.Now().Local().Year(),
					date.Month(),
					date.Day(),
					date.Hour(),
					date.Minute(),
					date.Second(),
					date.Nanosecond(),
					loc,
				)
			}
		case "div":
			{
				groups := child.FindAll("div", "class", "event-group")
				for _, g := range groups {
					var match broadcast.Match
					match.Competition = comp

					eventTitles := g.Find("ul", "class", "event").FindAll("strong")
					h, a, keep := extr(eventTitles)
					if !keep {
						continue
					}
					match.HomeTeam = h
					match.AwayTeam = a

					eventDetail := g.Find("p", "class", "event-detail").Text()
					foundTime := re.FindString(eventDetail)
					timeString := strings.Trim(foundTime, "()")
					var hours, mins int
					fmt.Sscanf(timeString, "%d:%d", &hours, &mins)

					curDateTime := time.Date(
						curDate.Year(),
						curDate.Month(),
						curDate.Day(),
						hours,
						mins,
						0,
						0,
						curDate.Location(),
					)
					match.Datetime = curDateTime

					Matches = append(Matches, broadcast.Broadcast{Match: match, Station: broadcast.SkySports})
				}
			}
		}
	}

	return Matches, nil
}