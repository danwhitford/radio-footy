package feeds

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type Match struct {
	Datetime    time.Time
	HomeTeam    string
	AwayTeam    string
	Competition string
}

func NewSantisedMatch(datetime time.Time, homeTeam, awayTeam, competition string) Match {
	return Match{
		datetime,
		mapTeamName(homeTeam),
		mapTeamName(awayTeam),
		mapCompName(competition),
	}
}

func (m Match) Title() string {
	if m.Competition == "NFL" {
		return fmt.Sprintf("%s @ %s", m.AwayTeam, m.HomeTeam)
	}
	return fmt.Sprintf("%s v %s", m.HomeTeam, m.AwayTeam)
}

func (m Match) Time() string {
	return m.Datetime.Format(timeLayout)
}

func (match Match) RollUpHash() string {
	return fmt.Sprintf("%s%v%s%s",
		match.Competition,
		match.Datetime.Format(time.DateOnly),
		match.HomeTeam,
		match.AwayTeam,
	)
}

func mapTeamName(name string) string {
	nameMapper := map[string]string{
		"IR Iran":                  "Iran",
		"Korea Republic":           "South Korea",
		"Milan":                    "AC Milan",
		"FC Bayern München":        "Bayern Munich",
		"Brighton and Hove Albion": "Brighton & Hove Albion",
		"Internazionale":           "Inter Milan",
		"Wolverhampton Wanderers":  "Wolves",
		"West Bromwich Albion":     "West Brom",
		"FC København":             "FC Copenhagen",
	}
	newName, prs := nameMapper[name]
	if prs {
		return newName
	} else {
		return name
	}
}

func mapCompName(competition string) string {
	replacements := map[*regexp.Regexp]string{
		regexp.MustCompile("Carabao Cup"):                    "EFL Cup",
		regexp.MustCompile("English Football League Trophy"): "EFL Cup",
		regexp.MustCompile("^EFL Trophy.*"):                  "EFL Cup",
		regexp.MustCompile("^EFL Cup.*"):                  "EFL Cup",
		regexp.MustCompile("[UEFA ]*Champions League.*"):     "Champions League",
		regexp.MustCompile("^Premier League.*"):              "Premier League",
		regexp.MustCompile("^FA Cup.*"):                      "FA Cup",
	}
	for old, new := range replacements {
		if old.MatchString(competition) {
			return new
		}
	}
	return competition
}

func (m Match) shouldSkip() bool {
	return strings.Contains(m.Competition, "Scottish") ||
		strings.Contains(m.Competition, "Women") ||
		strings.Contains(m.HomeTeam, "Scottish") ||
		strings.Contains(m.HomeTeam, "Women")
}
