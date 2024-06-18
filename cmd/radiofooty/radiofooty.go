package main

import (
	_ "embed"
	"html/template"
	"io"
	"log"
	"os"
	"slices"
	"time"

	"whitford.io/radiofooty/internal/broadcast"
	"whitford.io/radiofooty/internal/feeds"
)

//go:embed website/template.go.tmpl
var webTemplate string

//go:embed website/icalendar.go.tmpl
var calTemplate string

func main() {
	log.Println("Running wireless football")
	matches, err := feeds.GetMatches()
	if err != nil {
		log.Fatalln(err)
	}

	channelsSet := make(map[broadcast.Station]interface{})
	for _, day := range matches {
		for _, listing := range day.Matches {
			for _, station := range listing.Stations {
				channelsSet[station] = struct{}{}
			}
		}
	}
	channels := make([]broadcast.Station, 0)
	for channel, _ := range channelsSet {
		channels = append(channels, channel)
	}
	slices.SortFunc(channels, func(a, b broadcast.Station) int {
		return a.Rank - b.Rank
	})

	buildTime := time.Now().Format(time.RFC850)

	data := struct {
		MatchDays []broadcast.MatchDay
		UniqueChannels []broadcast.Station
		BuildTime string
	}{
		MatchDays: matches,
		UniqueChannels: channels,
		BuildTime: buildTime,
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

	writeIndex(data, "template.go.tmpl", webTemplate, idx)
	writeIndex(calData, "icalendar.go.tmpl", calTemplate, cal)
}

func writeIndex(data interface{}, templateName, contents string, writer io.Writer) {
	tmpl, err := template.New(templateName).Parse(contents)
	if err != nil {
		log.Fatalf("template parsing: %s", err)
	}
	err = tmpl.Execute(writer, data)
	if err != nil {
		log.Fatalf("template execution: %s", err)
	}
}
