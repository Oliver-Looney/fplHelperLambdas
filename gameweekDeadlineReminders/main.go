package main

import (
	"fmt"
	"time"
)

func getFPLGameweekData() (*fplBootstrapResponse, error) {
	c := NewClient(nil)
	url := "https://fantasy.premierleague.com/api/bootstrap-static/"
	response, err := c.NewRequest("GET", url)
	if err != nil {
		return nil, err
	}
	fplBootstrapResponse := &fplBootstrapResponse{}
	_, err = c.Do(response, fplBootstrapResponse)
	if err != nil {
		return nil, err
	}
	return fplBootstrapResponse, nil
}

func main() {
	fplGameweekData, err := getFPLGameweekData()
	if err != nil {
		println("ERROR")
	}
	for i := 0; i < len(fplGameweekData.Events); i++ {
		fmt.Println(fplGameweekData.Events[i])
	}
	return
}

type fplBootstrapResponse struct {
	Events []struct {
		ID                int       `json:"id"`
		Name              string    `json:"name"`
		DeadlineTime      time.Time `json:"deadline_time"`
		Finished          bool      `json:"finished"`
		DeadlineTimeEpoch int       `json:"deadline_time_epoch"`
	} `json:"events"`
}
