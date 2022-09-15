package main

import (
	"testing"
)

// to see logs use -> go test *.go -v

func TestIsSafeRowValid1(t *testing.T) {
	fplGameweekData, err := getFPLGameweekData()
	if err != nil {
		println("ERROR")
	}
	i := 0
	for !fplGameweekData.Events[i].IsNext {
		i++
	}
	t.Logf("\n%s\n", generateSMSContents(fplGameweekData.Events[i], fplGameweekData.Teams))
}
