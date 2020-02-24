package api

import (
	"net/http"
	"time"

	"github.com/go-ozzo/ozzo-routing/v2"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
)

// Version of API
const Version = "1.0.0"

// App is the application for API
type App struct {
	*commonApp.App
	Server		*http.Server
}

// New func is a constructor for the ApiApp
func New(cfg config.Configuration) *App {
	app := &App{
		App: commonApp.New(cfg),
		Server:	nil,
	}

	// build HTTP server
	server := &http.Server{
		Addr:		cfg.Server.HTTPListen,
		Handler:	app.buildHandler(),
	}
	app.Server = server

	return app
}

func (app *App) buildHandler() *routing.Router {
	router := routing.New()

	router.Get("/hello", func(c *routing.Context) error {
		return c.Write("Hello word!")
	})

	return 	router
}

// Run is func to run the ApiApp
func (app *App) Run() error {
	// start the HTTP server with graceful shutdown
	go routing.GracefulShutdown(app.Server, 10*time.Second, app.Logger.Infof)
	app.Logger.Infof("server %v is running at %v", Version, app.Server.Addr)
	if err := app.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}
