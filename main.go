package main

import (
	"os"

	"github.com/skamranahmed/markeventas/cmd"
	"github.com/skamranahmed/markeventas/pkg/log"
)

func main() {
	err := cmd.Run()
	if err != nil {
		log.Errorf("error starting the server, error: %v\n", err)
		os.Exit(1)
	}
}
