package api

import (

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/config"
)

// App is the application for API
type App struct {
	*commonApp.App
}

// New func is a constructor for the ApiApp
func New(cfg config.Configuration) *App {
	app := &App{
		commonApp.New(cfg),
	}

	return app
}

// Run is func to run the ApiApp
func (app *App) Run() error {
	return nil
}
