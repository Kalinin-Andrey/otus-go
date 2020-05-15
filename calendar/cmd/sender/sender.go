package main

import (
	"context"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"log"

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/sender"
)


func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Can not load the config, error: %v", err)
	}
	app := sender.New(context.Background(), commonApp.New(*cfg), *cfg)

	if err := app.Run(); err != nil {
		log.Fatalf("Error while application is running: %s", err.Error())
	}
}
