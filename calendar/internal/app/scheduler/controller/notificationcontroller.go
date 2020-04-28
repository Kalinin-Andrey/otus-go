package controller

import (
	"context"
	"encoding/json"
	"github.com/pkg/errors"

	"github.com/streadway/amqp"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/rabbitmq"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/notification"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/apperror/apperror"
)

// NotificationController struct
type NotificationController struct {
	ctx				context.Context
	Logger			log.ILogger
	queue			rabbitmq.QueueClient
	EventService	event.IService
}

// NotificationChannelSize const
const NotificationChannelSize = 100

// NewNotificationController returns a pointer to a new NotificationController obj
func NewNotificationController(ctx context.Context, eventService event.IService, logger	log.ILogger, queue rabbitmq.QueueClient) *NotificationController {
	return &NotificationController{
		ctx:			ctx,
		EventService:	eventService,
		Logger:			logger.With(ctx),
		queue:			queue,
	}
}

// RegisterQueueHandler registers handler for the queue and returns channel of Notification
func (c *NotificationController) RegisterQueueHandler() chan notification.Notification {
	ch := make(chan notification.Notification, NotificationChannelSize)

	c.queue.Handle(func(deliveryChan <- chan amqp.Delivery) {
		defer close(ch)
		for d := range deliveryChan {
			ch <- c.deliveryToNotification(d)
		}
		c.Logger.Debug("deliveryChan was closed, stop handler")
	}, 1)
	return ch
}

// deliveryToNotification convs delivery to notification
func (c *NotificationController) deliveryToNotification(d amqp.Delivery) notification.Notification {
	n := &notification.Notification{}

	err := json.Unmarshal(d.Body, n)
	if err != nil {
		c.Logger.Errorf("deliveryToNotification can not unmarshal a message body: %q", string(d.Body))
	}

	return *n
}

// Schedule gets notifications and send they in queue
func (c *NotificationController) Schedule() error {

	// Сделать выборку Event
	items, err := c.EventService.NotificationsList(c.ctx)
	if err != nil {
		if err == apperror.ErrNotFound {
			c.Logger.Info(err)
			return nil
		}
		err = errors.Wrapf(err, "Schedule receive error from EventService.ListForNotification: %q", err)
		c.Logger.Error(err)
		return err
	}
	// В цикле записать в очередь
	for _, n := range items{
		s, err := json.Marshal(n)
		if err != nil {
			return errors.Wrapf(err, "can not marshal a value: %v", n)
		}
		err = c.queue.Publish(s)
		if err != nil {
			return errors.Wrapf(err, "can not publish a value: %v", s)
		}
	}

	return nil
}

// Send is sends Notification obj
func (c *NotificationController) Send(n notification.Notification) {
	b, err := json.Marshal(n)
	if err != nil {
		c.Logger.Errorf("can not send, Marshal error: %v", err)
	}
	// "send" the message
	c.Logger.Debugf(" [ <- q] Sent notification to user: %s", string(b))
}


