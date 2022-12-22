package main

import (
	"bufio"
	"html/template"
	"os"
	"time"
	"unicode/utf8"

	"whitford.io/radiofooty/internal/feeds"
	"whitford.io/radiofooty/internal/interchange"
)

func main() {
	matches := feeds.GetMergedMatches()

	col := 0
	for _, matchDay := range matches {
		for _, match := range matchDay.Matches {
			if utf8.RuneCountInString(match.Title) > col {
				col = utf8.RuneCountInString(match.Title)
			}
			if utf8.RuneCountInString(match.Station) > col {
				col = utf8.RuneCountInString(match.Station)
			}
			if utf8.RuneCountInString(match.Competition) > col {
				col = utf8.RuneCountInString(match.Competition)
			}
		}
	}

	data := struct {
		MatchDays []interchange.MergedMatchDay
	}{
		MatchDays: matches,
	}

	events := feeds.MergedMatchDayToEventList(matches)
	dtstamp := time.Now().UTC().Format(interchange.CalTimeString)
	calData := struct {
		DtStamp string
		Events  []interchange.CalEvent
	}{
		DtStamp: dtstamp,
		Events:  events,
	}

	// Write index.html
	writeIndex(data)

	// Write iCalendar
	writeCal(calData)
}

func writeIndex(data interface{}) {
	template, err := template.ParseFiles("./internal/website/template.go.tmpl")
	if err != nil {
		panic(err)
	}
	f, _ := os.Create("index.html")
	defer f.Close()
	w := bufio.NewWriter(f)
	err = template.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func writeCal(data interface{}) {
	calTemplate, err := template.ParseFiles("./internal/website/icalendar.go.tmpl")
	if err != nil {
		panic(err)
	}
	f, _ := os.Create("icalendar.ics")
	defer f.Close()
	w := bufio.NewWriter(f)
	err = calTemplate.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}
