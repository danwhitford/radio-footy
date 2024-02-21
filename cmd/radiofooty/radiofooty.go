package main

import (
	"html/template"
	"io"
	"log"
	"os"

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

	calData := feeds.MatchDayToCalData(matches)

	idx, err := os.Create("index.html")
	if err != nil {
		panic(err)
	}
	cal, err := os.Create("icalendar.ics")
	if err != nil {
		panic(err)
	}

	writeIndex(data, "template.go.tmpl", "./internal/website/template.go.tmpl", idx)
	writeIndex(calData, "icalendar.go.tmpl", "./internal/website/icalendar.go.tmpl", cal)
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
