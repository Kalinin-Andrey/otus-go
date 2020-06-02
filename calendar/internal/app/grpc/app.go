package grpc

import (
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

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
	MetricServer	*http.Server
	Server			*grpc.Server
	Address			string
}

// Create a metrics registry.
var reg = prometheus.NewRegistry()

// Create some standard server metrics.
var grpcMetrics = grpc_prometheus.NewServerMetrics()

func init() {
	// Register standard server metrics.
	reg.MustRegister(grpcMetrics)
}

// New func is a constructor for the ApiApp
func New(commonApp *commonApp.App, cfg config.Configuration) *App {
	app := &App{
		App:	commonApp,
		Server:	grpc.NewServer(
			grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
			),
		Address: cfg.Server.GRPCListen,
	}

	c := controller.EventController{
		Service:		app.Domain.Event.Service,
		Logger:			app.Logger,
	}

	//reflection.Register(app.Server)	//	optional
	calendarpb.RegisterCalendarServer(app.Server, c)

	// Initialize all metrics.
	grpcMetrics.InitializeMetrics(app.Server)

	// Create a HTTP server for prometheus.
	app.MetricServer = &http.Server{
		Handler: promhttp.HandlerFor(reg, promhttp.HandlerOpts{}),
		Addr: cfg.Server.HTTPForPrometheuslisten,
	}

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

	// Start your http server for prometheus.
	go func() {
		app.Logger.Infof("HTTP metric server is runnig at %v", app.MetricServer.Addr)
		if err := app.MetricServer.ListenAndServe(); err != nil {
			app.Logger.Errorf("Unable to start a http metric server at %v", app.MetricServer.Addr)
		}
	}()

	app.Logger.Infof("grpc server %v is running at %v", Version, app.Address)

	if err := app.Server.Serve(l); err != nil && err != grpc.ErrServerStopped {
		return err
	}
	return nil
}

