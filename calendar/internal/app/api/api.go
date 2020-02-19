package api

import (

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/config"
)

type ApiApp struct {
	*commonApp.App
}

func New(cfg config.Configuration) *ApiApp {
	app := &ApiApp{
		commonApp.New(cfg),
	}

	return app
}


func (app *ApiApp) Run() error {
	return nil
}
