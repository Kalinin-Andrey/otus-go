package api

import (
	golog "log"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/db"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/infrastructure/repository"
)

// App struct is the common part of all applications
type App struct {
	Cfg    config.Configuration
	Logger log.ILogger
	DB     db.IDB
	Event  struct{
		//entity		event.IEvent
		Repository	event.IEventRepository
	}
}

// New func is a constructor for the App
func New(cfg config.Configuration) *App {
	logger, err := log.New(cfg.Log)
	if err != nil {
		panic(err)
	}

	db, err := db.New(cfg.DB, logger)
	if err != nil {
		panic(err)
	}

	app := &App{
		Cfg: cfg,
		Logger:	logger,
		DB:		db,
	}

	repository, err := repository.Get("event", app.Cfg.Repository.Type)
	if err != nil {
		golog.Fatalf("Can not get repository type %q for entity %q, error happened: %v", app.Cfg.Repository.Type, "event", err)
	}
	app.Event.Repository = repository

	return app
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}
