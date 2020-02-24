package main

import (
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/api"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"log"
)


func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalln("Can not load the config")
	}
	app := api.New(*cfg)

	if err := app.Run(); err != nil {
		log.Fatalf("Error while application is running: %s", err.Error())
	}
}

