package main

import (
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"log"

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/rest"
)


func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Can not load the config, error: %v", err)
	}
	app := rest.New(commonApp.New(*cfg), *cfg)

	if err := app.Run(); err != nil {
		log.Fatalf("Error while application is running: %s", err.Error())
	}
}

