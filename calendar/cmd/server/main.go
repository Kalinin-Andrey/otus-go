package main

import (
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/api"
	"log"
)


func main() {
	err, cfg := config.Get()
	if err != nil {
		log.Fatalln("Can not load the config")
	}
	app := api.New(*cfg)

	if err := app.Run(); err != nil {
		log.Fatalf("Error while application is running: %s", err.Error())
	}
}

