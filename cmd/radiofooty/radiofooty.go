package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"whitford.io/radiofooty/internal/feeds"
)

func main() {
	matches, err := feeds.GetMergedMatches()
	if err != nil {
		log.Fatalln(err)
	}

	data := struct {
		MatchDays []feeds.MergedMatchDay
	}{
		MatchDays: matches,
	}

	events := feeds.MergedMatchDayToEventList(matches)
	dtstamp := time.Now().UTC().Format(feeds.CalTimeString)
	calData := struct {
		DtStamp string
		Events  []feeds.CalEvent
	}{
		DtStamp: dtstamp,
		Events:  events,
	}

	// Write index.html
	f, err := os.Create("index.html")
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer f.Close()
	writeIndex(data, "./internal/website/template.go.tmpl", f)

	// Write iCalendar
	fcal, err := os.Create("icalendar.ics")
	if err != nil {
		log.Fatalf("error creating file: %v", err)
	}
	defer fcal.Close()
	writeCal(calData, "./internal/website/icalendar.go.tmpl", fcal)
}

func writeIndex(data interface{}, templatePath string, writer io.Writer) {
	funcs := template.FuncMap{
		"join": strings.Join,
		"rfc3339": func(t time.Time) string {
			return t.Format(time.DateOnly)
		},
		"gamehash": func(m feeds.MergedMatch) string {			
			s := fmt.Sprintf("%s%s", m.Date, m.Title)
			s = strings.ReplaceAll(s, " ", "")
			s = strings.ReplaceAll(s, ",", "")
			return s
		},
	}
	tmpl, err := template.New("template.go.tmpl").Funcs(funcs).ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	err = tmpl.Execute(writer, data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}

func writeCal(data interface{}, templatePath string, w io.Writer) {
	funcs := template.FuncMap{
		"join": strings.Join,
	}
	calTemplate, err := template.New("icalendar.go.tmpl").Funcs(funcs).ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}

	err = calTemplate.Execute(w, data)
	if err != nil {
		panic(err)
	}
}
