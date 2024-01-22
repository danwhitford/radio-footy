package feeds

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/anaskhan96/soup"
	"whitford.io/radiofooty/internal/urlgetter"
)

const nflUrl string = "https://www.skysports.com/watch/nfl-on-sky"

func getNflOnSky(getter urlgetter.UrlGetter) ([]MergedMatch, error) {
	html, err := getter.GetUrl(nflUrl)
	if err != nil {
		return nil, err
	}
	return nflPageToMergedMatches(string(html))
}

var dateRe *regexp.Regexp = regexp.MustCompile(`(\d+)(st|nd|rd|th)`)

func nflPageToMergedMatches(html string) ([]MergedMatch, error) {
	mergedMatches := make([]MergedMatch, 0)

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
		return mergedMatches, nil
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
					return mergedMatches, fmt.Errorf("could not parse %s: %s", cleanstring, err)
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
					var mergedMatch MergedMatch
					mergedMatch.Competition = "NFL"
					mergedMatch.Stations = []string{"Sky Sports"}

					eventTitles := g.Find("ul", "class", "event").FindAll("strong")
					mergedMatch.Title = eventTitles[0].Text() + " @ " + eventTitles[1].Text()

					eventDetail := g.Find("p", "class", "event-detail").Text()
					foundTime := re.FindString(eventDetail)
					timeString := strings.Trim(foundTime, "()")
					mergedMatch.Time = timeString
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
					mergedMatch.Datetime = curDateTime.Format(time.RFC3339)
					mergedMatch.Date = curDateTime.Format(niceDate)

					mergedMatches = append(mergedMatches, mergedMatch)
				}
			}
		}
	}

	return mergedMatches, nil
}
