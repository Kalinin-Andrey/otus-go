package app

import (
	golog "log"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/dbx"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	dbrep "github.com/Kalinin-Andrey/otus-go/calendar/internal/infrastructure/repository/db"
)

// IApp app interface
type IApp interface {
	// Run is func to run the App
	Run() error
}

// App struct is the common part of all applications
type App struct {
	Cfg					config.Configuration
	Logger				log.ILogger
	DB					db.IDB
	Domain				Domain
}

// Domain is a Domain Layer Entry Point
type Domain struct {
	Event struct {
		Repository event.IRepository
		Service    event.IService
	}
}

const (
	// EntityNameEvent event entity name constant
	EntityNameEvent			= "event"
)

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
		Cfg:    cfg,
		Logger: logger,
		DB:     db,
	}
	var ok bool

	app.Domain.Event.Repository, ok	= app.getDBRepo(EntityNameEvent).(event.IRepository)
	if !ok {
		golog.Fatalf("Can not cast DB repository for entity %q to %v.IRepository. Repo: %v", EntityNameEvent, EntityNameEvent, app.getDBRepo(EntityNameEvent))
	}
	app.Domain.Event.Service = event.NewService(app.Domain.Event.Repository, app.Logger)

	return app
}

// Run is func to run the App
func (app *App) Run() error {
	return nil
}

func (app *App) getDBRepo(entityName string) (repo dbrep.IRepository) {
	var err error

	if repo, err = dbrep.GetRepository(app.DB, app.Logger, entityName); err != nil {
		golog.Fatalf("Can not get db repository for entity %q, error happened: %v", entityName, err)
	}
	return repo
}
