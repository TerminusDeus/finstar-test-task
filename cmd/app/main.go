package main

import (
	"log"

	"finstar-test-task/internal/app"

	"finstar-test-task/config"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// Run
	app.Run(cfg)
}
