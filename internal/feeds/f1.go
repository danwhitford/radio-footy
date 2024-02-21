package feeds

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"whitford.io/radiofooty/internal/urlgetter"
)

const f1Url string = "https://www.skysports.com/watch/f1-on-sky"

func getF1OnSky(getter urlgetter.UrlGetter) ([]Broadcast, error) {
	html, err := getter.GetUrl(f1Url)
	if err != nil {
		return nil, err
	}
	return f1PageToMatches(string(html))
}

func f1PageToMatches(html string) ([]Broadcast, error) {
	Matches := make([]Broadcast, 0)

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
					var match Match
					match.Competition = "F1"

					eventTitles := g.Find("ul", "class", "event").FindAll("strong")
					raceTitle := eventTitles[0].Text()
					if !strings.HasSuffix(raceTitle, "Race") {
						continue
					}
					raceTitle, _ = strings.CutSuffix(raceTitle, " - Race")
					match.HomeTeam = raceTitle

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

					Matches = append(Matches, Broadcast{match, SkySports})
				}
			}
		}
	}

	return Matches, nil
}
