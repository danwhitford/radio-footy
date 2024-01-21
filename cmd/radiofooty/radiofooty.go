package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"regexp"
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

	if len(os.Args) < 2 {
		log.Fatalf("Need to supply a target. website or calendar.")
	}
	target := os.Args[1]
	switch target {
	case "website":
		writeIndex(data, "template.go.tmpl", "./internal/website/template.go.tmpl", os.Stdout)
	case "calendar":
		writeIndex(calData, "icalendar.go.tmpl", "./internal/website/icalendar.go.tmpl", os.Stdout)
	default:
		log.Fatalf("Target not recognised: %s\n", target)
	}

}

func writeIndex(data interface{}, templateName, templatePath string, writer io.Writer) {
	funcs := template.FuncMap{
		"join": strings.Join,
		"rfc3339": func(t time.Time) string {
			return t.Format(time.DateOnly)
		},
		"gamehash": func(m feeds.MergedMatch) string {
			s := fmt.Sprintf("%s%s", m.Date, m.Title)
			r := regexp.MustCompile("[^0-9a-zA-Z]")
			s = r.ReplaceAllString(s, "")
			return s
		},
	}
	tmpl, err := template.New(templateName).Funcs(funcs).ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	err = tmpl.Execute(writer, data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}
