package grpc

import (
	"google.golang.org/grpc"
	"log"
	"net"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/grpc/calendarpb"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/grpc/controller"
)

// Version of API
const Version = "1.0.0"

// App is the application for API
type App struct {
	*commonApp.App
	Server		*grpc.Server
	Address		string
}

// New func is a constructor for the ApiApp
func New(commonApp *commonApp.App, cfg config.Configuration) *App {
	app := &App{
		App:	commonApp,
		Server:	grpc.NewServer(),
		Address: cfg.Server.GRPCListen,
	}

	c := controller.EventController{
		Service:		app.Domain.Event.Service,
		Logger:			app.Logger,
	}

	//reflection.Register(app.Server)	//	optional
	calendarpb.RegisterCalendarServer(app.Server, c)

	return app
}

// Run is func to run the ApiApp
func (app *App) Run() error {
	defer func() {
		if err := app.DB.DB().Close(); err != nil {
			app.Logger.Error(err)
		}

		err := app.Logger.Sync()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	l, err := net.Listen("tcp", app.Address)
	if err != nil {
		log.Fatalf("failed to listen %q, err: %q", app.Address, err)
	}


	app.Logger.Infof("grpc server %v is running at %v", Version, app.Address)

	if err := app.Server.Serve(l); err != nil && err != grpc.ErrServerStopped {
		return err
	}
	return nil
}

