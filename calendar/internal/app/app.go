package api

import (
	"log"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/infrastructure/repository"
)

// App struct is the common part of all applications
type App struct {
	Cfg	config.Configuration
	Event	struct{
		//entity		event.IEvent
		Repository	event.IEventRepository
	}
}

// New func is a constructor for the App
func New(cfg config.Configuration) *App {
	app := &App{
		Cfg: cfg,
	}

	repository, err := repository.Get("event", app.Cfg.Repository.Type)
	if err != nil {
		log.Fatalf("Can not get repository type %q for entity %q, error happened: %v", app.Cfg.Repository.Type, "event", err)
	}
	app.Event.Repository = repository

	return app
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}
