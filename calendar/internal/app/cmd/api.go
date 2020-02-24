package cmd

import (

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
)

// App is the application for CLI app
type App struct {
	*commonApp.App
}

// New func is a constructor for the App
func New(cfg config.Configuration) *App {
	app := &App{
		commonApp.New(cfg),
	}

	return app
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}
