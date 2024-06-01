package channel

import (
	"encoding/json"
	"fmt"
	"time"

	"whitford.io/radiofooty/internal/broadcast"
)

type ManualGetter struct{}

const weightliftingJsonString string = `{
  "events": [
    {
      "venue": "South Paris Arena 6",
      "date": "Wednesday",
      "day": "07",
      "month": "August",
      "schedule": [
        {
          "time": "14:00",
          "event": "Men's 61kg"
        },
        {
          "time": "18:30",
          "event": "Women's 49kg"
        }
      ]
    },
    {
      "venue": "South Paris Arena 6",
      "date": "Thursday",
      "day": "08",
      "month": "August",
      "schedule": [
        {
          "time": "14:00",
          "event": "Women's 59kg"
        },
        {
          "time": "18:30",
          "event": "Men's 73kg"
        }
      ]
    },
    {
      "venue": "South Paris Arena 6",
      "date": "Friday",
      "day": "09",
      "month": "August",
      "schedule": [
        {
          "time": "14:00",
          "event": "Men's 89kg"
        },
        {
          "time": "18:30",
          "event": "Women's 71kg"
        }
      ]
    },
    {
      "venue": "South Paris Arena 6",
      "date": "Saturday",
      "day": "10",
      "month": "August",
      "schedule": [
        {
          "time": "10:30",
          "event": "Men's 102kg"
        },
        {
          "time": "15:00",
          "event": "Women's 81kg"
        },
        {
          "time": "19:30",
          "event": "Men's +102kg"
        }
      ]
    },
    {
      "venue": "South Paris Arena 6",
      "date": "Sunday",
      "day": "11",
      "month": "August",
      "schedule": [
        {
          "time": "10:30",
          "event": "Women's +81kg"
        }
      ]
    }
  ]
}`

type Event struct {
	Time  string `json:"time"`
	Event string `json:"event"`
}

type Schedule struct {
	Venue    string  `json:"venue"`
	Date     string  `json:"date"`
	Day      string  `json:"day"`
	Month    string  `json:"month"`
	Schedule []Event `json:"schedule"`
}

type Events struct {
	Events []Schedule `json:"events"`
}

const dtformat string = "Monday 02 January 2006 15:04"

const OlympicWLTitle string = "Paris 2024 Olympic Weightlifting"

func (mg ManualGetter) GetMatches() ([]broadcast.Broadcast, error) {
	var weightlifting Events
	err := json.Unmarshal([]byte(weightliftingJsonString), &weightlifting)
	if err != nil {
		return []broadcast.Broadcast{}, fmt.Errorf("failed to unmarshal weightlifting. %v", err)
	}

	var broadcasts []broadcast.Broadcast
	for _, event := range weightlifting.Events {
		for _, sched := range event.Schedule {
			dt, err := time.Parse(
				dtformat,
				fmt.Sprintf("%s %s %s 2024 %s", event.Date, event.Day, event.Month, sched.Time),
			)
			if err != nil {
				return broadcasts, fmt.Errorf("couldn't parse datetime. %v", err)
			}
			m := broadcast.Match{
				Competition: OlympicWLTitle,
				HomeTeam:    sched.Event,
				Datetime:    dt,
			}
			b := broadcast.Broadcast{
				Match:   m,
				Station: broadcast.Station{Name: "BBC", Rank: 1},
			}
			broadcasts = append(broadcasts, b)
		}
	}

	return broadcasts, nil
}
