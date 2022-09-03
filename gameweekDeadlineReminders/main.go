package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"os"
	"strconv"
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

func generateSMSContents(gameweek gameweekData) string {
	return gameweek.Name +
		" Deadline is in " +
		strconv.FormatInt(getDaysFromEpochTime(gameweek), 10) +
		" days and " + strconv.FormatInt(getHoursFromEpochTime(gameweek), 10) +
		" hours\n" +
		"Deadline: " +
		string(gameweek.DeadlineTime.Weekday()) +
		" " +
		string(gameweek.DeadlineTime.Hour()) +
		":" +
		string(gameweek.DeadlineTime.Minute())
}

func getDaysFromEpochTime(gameweek gameweekData) int64 {
	return (gameweek.DeadlineTimeEpoch - time.Now().Unix()) / 86400
}

func getHoursFromEpochTime(gameweek gameweekData) int64 {
	return ((gameweek.DeadlineTimeEpoch - time.Now().Unix()) % 86400) / 3600
}

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) {
	fplGameweekData, err := getFPLGameweekData()
	if err != nil {
		println("ERROR")
	}
	i := 0
	for !fplGameweekData.Events[i].IsNext {
		i++
	}
	if getDaysFromEpochTime(fplGameweekData.Events[i]) > 3 {
		return
	}
	sess := session.Must(session.NewSession())

	svc := sns.New(sess)

	params := &sns.PublishInput{
		Message:     aws.String(generateSMSContents(fplGameweekData.Events[i])),
		PhoneNumber: aws.String(os.Getenv("MY_PHONE_NUMBER")),
	}
	resp, err := svc.Publish(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println(resp)
	return
}

func main() {
	lambda.Start(HandleRequest)
}

type gameweekData struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	DeadlineTime      time.Time `json:"deadline_time"`
	Finished          bool      `json:"finished"`
	IsNext            bool      `json:"is_next"`
	DeadlineTimeEpoch int64     `json:"deadline_time_epoch"`
}

type fplBootstrapResponse struct {
	Events []gameweekData `json:"events"`
}
