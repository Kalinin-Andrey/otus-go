package cmd

import (

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/config"
)

type CmdApp struct {
	*commonApp.App
}

func New(cfg config.Configuration) *CmdApp {
	app := &CmdApp{
		commonApp.New(cfg),
	}

	return app
}


func (app *CmdApp) Run() error {
	return nil
}
