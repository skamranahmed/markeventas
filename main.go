package main

import (
	"os"

	"github.com/skamranahmed/twitter-create-gcal-event-api/cmd"
	"github.com/skamranahmed/twitter-create-gcal-event-api/pkg/log"
)

func main() {
	log.Info("Twitter Create Gcal Event API Test")
	err := cmd.Run()
	if err != nil {
		log.Errorf("error starting the server, error: %v\n", err)
		os.Exit(1)
	}
}
