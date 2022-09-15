package main

import "time"

type gameweekData struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	DeadlineTime      time.Time `json:"deadline_time"`
	Finished          bool      `json:"finished"`
	IsNext            bool      `json:"is_next"`
	DeadlineTimeEpoch int64     `json:"deadline_time_epoch"`
}

type teamData struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

type fplBootstrapResponse struct {
	Events []gameweekData `json:"events"`
	Teams  []teamData     `json:"teams"`
}

type gameweekFixtures []struct {
	Event int `json:"event"`
	TeamA int `json:"team_a"`
	TeamH int `json:"team_h"`
}

type MyEvent struct {
	Name string `json:"name"`
}
