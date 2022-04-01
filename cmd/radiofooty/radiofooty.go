package main

import (
	// "encoding/json"
	// "fmt"
	"bufio"
	"html/template"
	"os"
	"strings"

	"whitford.io/radiofooty/internal/feeds"
	"whitford.io/radiofooty/internal/interchange"
)

func main() {
	matches := feeds.GetMergedMatches()

	template, er := template.ParseFiles("./internal/website/template.gohtml")
	if er != nil {
		panic(er)
	}

	col := 0
	for _, matchDay := range matches {
		for _, match := range matchDay.Matches {
			if len(match.Title) > col {
				col = len(match.Title)
			}
			if len(match.Station) > col {
				col = len(match.Station)
			}
		}
	}

	data := struct {
		MatchDays   []interchange.MergedMatchDay
		Pad         func(string, int) string
		Repeat      func(string, int) string
		RepeatBlind func(int, int) string
		Col         int
	}{
		MatchDays: matches,
		Pad: func(s string, n int) string {
			l := len(s)
			p := n - l
			return s + strings.Repeat("\u00A0", p)
		},
		Repeat: func(s string, i int) string {
			return strings.Repeat(s, i)
		},
		RepeatBlind: func(i, j int) string {
			return strings.Repeat("\u00A0", j-i)
		},
		Col: col,
	}

	f, _ := os.Create("index.html")
	defer f.Close()
	w := bufio.NewWriter(f)
	err := template.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
}
