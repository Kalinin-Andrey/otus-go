package main

import (
	"log"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/grpc"
)


func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatalln("Can not load the config")
	}
	app := grpc.New(commonApp.New(*cfg), *cfg)

	if err := app.Run(); err != nil {
		log.Fatalf("Error while application is running: %s", err.Error())
	}
}
