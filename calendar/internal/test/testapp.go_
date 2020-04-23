package test

import (
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/test/mock"
)

// NewCommonApp func is a constructor for the App
func NewCommonApp(cfg config.Configuration) *commonApp.App {
	logger, err := log.New(cfg.Log)
	if err != nil {
		panic(err)
	}

	app := &commonApp.App{
		Cfg:    cfg,
		Logger: logger,
		DB:     nil,
	}

	app.Domain.Event.Repository = &mock.EventRepository{}
	app.Domain.Event.Service = event.NewService(app.Domain.Event.Repository, app.Logger)

	return app
}

