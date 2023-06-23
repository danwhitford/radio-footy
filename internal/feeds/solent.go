package feeds

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"whitford.io/radiofooty/internal/filecacher"
)

func getSolentMatches() ([]MergedMatch, error) {
	var matches = []MergedMatch{}
	baseUrls := []string{
		"https://rms.api.bbc.co.uk/v2/experience/inline/schedules/bbc_radio_solent/",
	}
	urls := []string{}
	urlTime := "2006-01-02"
	start := time.Now()
	for i := 0; i < 8; i++ {
		t := start.AddDate(0, 0, i)
		for _, baseUrl := range baseUrls {
			urls = append(urls, baseUrl+t.Format(urlTime))
		}
	}

	for _, url := range urls {
		body, err := filecacher.GetUrl(url, filecacher.HttpGetter{})
		if err != nil {
			return nil, err
		}

		merged, err := getSolentDay(string(body))
		if err != nil {
			return nil, err
		}
		matches = append(matches, merged...)
	}

	return matches, nil
}

func getSolentDay(html string) ([]MergedMatch, error) {
	doc, err := htmlquery.Parse(strings.NewReader(html))
	if err != nil {
		return nil, err
	}
	scripts := htmlquery.Find(doc, "//script")
	for _, sc := range scripts {
		if sc.FirstChild != nil {
			scriptContent := sc.FirstChild.Data
			if strings.HasPrefix(strings.TrimSpace(scriptContent), "window.__PRELOADED_STATE__") {
				actualJson := strings.TrimSuffix(
					strings.TrimPrefix(strings.TrimSpace(scriptContent), "window.__PRELOADED_STATE__ = "),
					";",
				)
				var foo struct {
					Modules BBCFeed `json:"modules"`
				}
				err := json.Unmarshal([]byte(actualJson), &foo)
				if err != nil {
					return nil, err
				}
				return solentDayToMergedMatch(foo.Modules), nil
			}
		}
	}

	return nil, fmt.Errorf("script tag for json not found")
}

func solentDayToMergedMatch(bbcFeed BBCFeed) []MergedMatch {
	matches := make([]MergedMatch, 0)

	loc, err := time.LoadLocation("Europe/London")
	if err != nil {
		log.Fatalf("error loading location: %v", err)
	}
	longFormat := "2006-01-02T15:04:05Z"

	for _, data := range bbcFeed.Data {
		for _, prog := range data.Data {
			if strings.HasPrefix(prog.Title.Secondary, "Women") {
				continue
			}
			if isLocalCricket(prog.Title) {
				start, err := time.Parse(longFormat, prog.Start)
				if err != nil {
					panic(err)
				}
				start = start.In(loc)
				clock := start.Format(timeLayout)
				date := start.Format(niceDate)

				rg := regexp.MustCompile(`[(].+[)]`)
				title := rg.ReplaceAllString(prog.Title.Secondary, "")
				title = strings.TrimSpace(title)				
				comp := strings.SplitAfter(prog.Synopses.Short, "cricket")[0]

				m := MergedMatch{
					Time:        clock,
					Date:        date,
					Stations:    []string{"BBC Radio Solent"},
					Datetime:    start.Format(time.RFC3339),
					Title:       title,
					Competition: comp,
				}

				//debug
				if !strings.Contains(title, " v ") {
					log.Fatalf("oh no %+v\nbecame %+v\n", prog, m)
				}

				matches = append(matches, m)
			}
		}
	}

	return matches
}

func isLocalCricket(title BBCTitles) bool {
	return title.Primary == "Summer Sport"
}
