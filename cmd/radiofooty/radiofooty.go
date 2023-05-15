package main

import (
	"bufio"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"whitford.io/radiofooty/internal/feeds"
	"whitford.io/radiofooty/internal/interchange"
)

func main() {
	matches := feeds.GetMergedMatches()

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
	f, err := os.Create("index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	writeIndex(data, "./internal/website/template.go.tmpl", f)

	// Write discordian
	writeDiscordian(data)

	// Write iCalendar
	writeCal(calData)
}

func writeDiscordian(data struct {
	MatchDays []interchange.MergedMatchDay
}) {
	for i := range data.MatchDays {
		dts := data.MatchDays[i].Matches[0].Datetime
		dt, _ := time.Parse(time.RFC3339, dts)
		y := strconv.Itoa(dt.Year())
		m := strconv.Itoa(int(dt.Month()))
		d := strconv.Itoa(dt.Day())
		cmd := exec.Command("ddate", d, m, y)
		cmdout, _ := cmd.Output()

		data.MatchDays[i].NiceDate = string(cmdout)
	}

	template, err := template.ParseFiles("./internal/website/template.go.tmpl")
	if err != nil {
		panic(err)
	}
	f, _ := os.Create("discordian.html")
	defer f.Close()
	w := bufio.NewWriter(f)
	err = template.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}

func writeIndex(data interface{}, templatePath string, writer io.Writer) {
	template, err := template.ParseFiles(templatePath)
	if err != nil {
		panic(err)
	}
	err = template.Execute(writer, data)
	if err != nil {
		panic(err)
	}
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
