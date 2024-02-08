package main

import (
	"html/template"
	"io"
	"log"
	"os"
	"time"

	"whitford.io/radiofooty/internal/feeds"
)

func main() {
	log.Println("Running wireless football")
	matches, err := feeds.GetMatches()
	if err != nil {
		log.Fatalln(err)
	}

	data := struct {
		MatchDays []feeds.MatchDay
	}{
		MatchDays: matches,
	}

	events := feeds.MatchDayToEventList(matches)
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
	tmpl, err := template.New(templateName).ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	err = tmpl.Execute(writer, data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}
