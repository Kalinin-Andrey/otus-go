package scheduler

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/rabbitmq"

	commonApp "github.com/Kalinin-Andrey/otus-go/calendar/internal/app"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/scheduler/controller"
)

// App is the application for CLI app
type App struct {
	*commonApp.App
	Ctx                    context.Context
	cancel                 context.CancelFunc
	Queue                  rabbitmq.QueueClient
	NotificationController *controller.NotificationController
}

// New func is a constructor for the App
func New(ctx context.Context, commonApp *commonApp.App, cfg config.Configuration) *App {
	ctx, cancel := context.WithCancel(ctx)
	queue, err := rabbitmq.NewClient(ctx, commonApp.Logger, cfg.Queue.RabbitMQ, rabbitmq.TypePublisher)
	if err != nil {
		panic(errors.Wrap(err, "can not connect to queue"))
	}

	app := &App{
		App:                    commonApp,
		Ctx:                    ctx,
		cancel:                 cancel,
		Queue:                  queue,
		NotificationController: controller.NewNotificationController(ctx, commonApp.Domain.Event.Service, commonApp.Logger, queue),
	}

	return app
}

// GracefulShutdown func
func (app *App) GracefulShutdown() {
	defer app.cancel()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
}

// Run is func to run the App
func (app *App) Run() error {
	go app.GracefulShutdown()

	// слушать таймер и запускать schedule()
	ticker := time.NewTicker(time.Minute)
OUTER:
	for {

		select {
		case <- app.Ctx.Done():
			ticker.Stop()
			break OUTER
		case <- ticker.C:
			if err := app.NotificationController.Schedule(); err != nil {
				app.Logger.Errorf("application run error: %v", err)
				return err
			}
		}

	}

	return nil
}


