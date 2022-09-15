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
	"strings"
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

func getFPLGameweekFixturesData(gameweek int) (*gameweekFixtures, error) {
	c := NewClient(nil)
	url := fmt.Sprintf("https://fantasy.premierleague.com/api/fixtures?event=%d", gameweek)
	response, err := c.NewRequest("GET", url)
	if err != nil {
		return nil, err
	}
	gameweekFixtures := &gameweekFixtures{}
	_, err = c.Do(response, gameweekFixtures)
	if err != nil {
		return nil, err
	}
	return gameweekFixtures, nil
}

func generateSMSContents(gameweek gameweekData, teamData []teamData) string {
	return gameweek.Name +
		" Deadline is in " +
		strconv.FormatInt(getDaysFromEpochTime(gameweek), 10) +
		" days and " + strconv.FormatInt(getHoursFromEpochTime(gameweek), 10) +
		" hours" +
		checkGameWeeksFixturesAndSummarize(gameweek.ID, teamData) +
		checkGameWeeksFixturesAndSummarize(gameweek.ID+1, teamData)
}

func checkGameWeeksFixturesAndSummarize(gameweekID int, teamData []teamData) string {
	fixturesCurr, err := getFPLGameweekFixturesData(gameweekID)
	if err != nil {
		fmt.Println(err)
	}
	fixturesAsMapCurr := fixturesListToMap(*fixturesCurr)
	return checkForTeamsNotPlayingOnce(fixturesAsMapCurr, gameweekID, teamData, "teams not playing in gameweek", testForZeroGames) + checkForTeamsNotPlayingOnce(fixturesAsMapCurr, gameweekID, teamData, "teams are playing multiple times in gameweek", testForMultipleGames)
}

func getDaysFromEpochTime(gameweek gameweekData) int64 {
	return (gameweek.DeadlineTimeEpoch - time.Now().Unix()) / 86400
}

func getHoursFromEpochTime(gameweek gameweekData) int64 {
	return ((gameweek.DeadlineTimeEpoch - time.Now().Unix()) % 86400) / 3600
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
		Message:     aws.String(generateSMSContents(fplGameweekData.Events[i], fplGameweekData.Teams)),
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

func fixturesListToMap(fixtures gameweekFixtures) map[int]int {
	teamsNumberOfGames := map[int]int{1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 0, 12: 0, 13: 0, 14: 0, 15: 0, 16: 0, 17: 0, 18: 0, 19: 0, 20: 0}
	for i := 0; i < len(fixtures); i++ {
		teamsNumberOfGames[fixtures[i].TeamH]++
		teamsNumberOfGames[fixtures[i].TeamA]++
	}
	return teamsNumberOfGames
}

type fn func(map[int]int, int, bool, []string, []teamData) (bool, []string)

func checkForTeamsNotPlayingOnce(mapOfTeamsGames map[int]int, gameweek int, teamData []teamData, message string, testCase fn) string {
	if gameweek > 38 {
		return ""
	}
	teamsToNote := make([]string, 0)
	flag := false
	for i := 1; i <= 20; i++ {
		flag, teamsToNote = testCase(mapOfTeamsGames, i, flag, teamsToNote, teamData)
	}
	if flag {
		return fmt.Sprintf("\n%d %s %d: %s", len(teamsToNote), message, gameweek, strings.Join(teamsToNote[:], ", "))
	} else {
		return ""
	}
}

func testForMultipleGames(mapOfTeamsGames map[int]int, i int, flag bool, teamsPlayingMultiple []string, teamData []teamData) (bool, []string) {
	if mapOfTeamsGames[i] > 1 {
		flag = true
		teamsPlayingMultiple = append(teamsPlayingMultiple, teamData[i-1].Name)
	}
	return flag, teamsPlayingMultiple
}

func testForZeroGames(mapOfTeamsGames map[int]int, i int, flag bool, teamsNotPlaying []string, teamData []teamData) (bool, []string) {
	if mapOfTeamsGames[i] < 1 {
		flag = true
		teamsNotPlaying = append(teamsNotPlaying, teamData[i-1].Name)
	}
	return flag, teamsNotPlaying
}

func main() {
	lambda.Start(HandleRequest)
}
