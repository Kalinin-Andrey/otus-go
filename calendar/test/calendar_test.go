package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/grpc/calendarpb"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/app/grpc/controller"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	golog "log"
	"os"
	"time"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"
	"github.com/cucumber/godog"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/pkg/rabbitmq"

	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/event"
	"github.com/Kalinin-Andrey/otus-go/calendar/internal/domain/notification"
)

const (
	NotificationChannelSize		= 100
	userID						= 1
)

func NewGRPCClient(ctx context.Context, address string) calendarpb.CalendarClient {
	conn, err := grpc.DialContext(ctx, address, grpc.WithInsecure())
	if err != nil {
		golog.Fatalf("Cannot connect: %v", err)
	}
	golog.Printf("GRPC client starts at %v\n", address)
	return calendarpb.NewCalendarClient(conn)
}

var amqpDSN		= os.Getenv("TESTS_AMQP_DSN")
var GRPCAddress	= os.Getenv("TESTS_GRPC_ADDRESS")

func init() {
	if amqpDSN == "" {
		//amqpDSN = "amqp://guest:guest@queue:5672"
		amqpDSN = "amqp://guest:guest@localhost:5672"
	}

	if GRPCAddress == "" {
		//GRPCAddress = "grpcapi:8888"
		GRPCAddress = "localhost:8882"
	}
}


type calendarTest struct {
	ctx				context.Context
	cancel			context.CancelFunc
	cfg				config.Configuration
	logger			log.Logger
	GRPCClient		calendarpb.CalendarClient
	queue			rabbitmq.Client
	queueCh			chan notification.Notification
	testResponses	testResponses
}

type testResponses struct {
	recponseEvent		*calendarpb.ResponseEvent
	recponseEvents		*calendarpb.ResponseEvents
	error				error
}

func newCalendarTest(ctx context.Context) *calendarTest {
	os.Chdir("../")
	cfg, err := config.Get()
	if err != nil {
		golog.Fatalf("Can not load the config, error: %v", err)
	}

	logger, err := log.New(cfg.Log)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(ctx)
	cfg.Queue.RabbitMQUserNotification.URI = amqpDSN
	queue, err := rabbitmq.NewClient(ctx, logger, cfg.Queue.RabbitMQUserNotification, rabbitmq.TypePublisher)

	ct := &calendarTest{
		ctx:		ctx,
		cancel:		cancel,
		cfg:        *cfg,
		logger:     *logger,
		GRPCClient: NewGRPCClient(ctx, GRPCAddress),
		queue:      *queue,
		queueCh:    nil,
	}
	ct.RegisterQueueHandler()

	return ct
}

// RegisterQueueHandler registers handler for the queue and returns channel of Notification
func (c *calendarTest) RegisterQueueHandler() error {
	c.queueCh = make(chan notification.Notification, NotificationChannelSize)

	return c.queue.Handle(func(deliveryChan <- chan amqp.Delivery) {
		defer close(c.queueCh)
		for d := range deliveryChan {
			c.queueCh <- c.deliveryToNotification(d)
		}
		c.logger.Debug("deliveryChan was closed, stop handler")
	}, 1)
}

// deliveryToNotification convs delivery to notification
func (c *calendarTest) deliveryToNotification(d amqp.Delivery) notification.Notification {
	n := &notification.Notification{}

	err := json.Unmarshal(d.Body, n)
	if err != nil {
		c.logger.Errorf("deliveryToNotification can not unmarshal a message body: %q", string(d.Body))
	}

	return *n
}




func (c *calendarTest) start() {
	fmt.Println("Start test!")
}

func (c *calendarTest) stop() {
	c.cancel()
	c.queue.Close()
}


func (c *calendarTest) panicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}




func (c *calendarTest) iCreateEventWithUserIDTitleDescription(userID int, title, description string) error {
	event := event.Event{
		UserID:       uint(userID),
		Title:        title,
		Description:  &description,
		Time:         time.Now().AddDate(1, 0, 0),
	}
	e, err := controller.EventToEventProto(event)
	if err != nil {
		return err
	}

	response, err := c.GRPCClient.EventCreate(c.ctx, e)
	if err != nil {
		return err
	}
	c.testResponses.recponseEvent = response

	return nil
}


func (c *calendarTest) iReceiveStatusIsOK() error {
	if !c.testResponses.recponseEvent.Status.OK {
		return errors.Errorf("Expected status is OK, have received status: %v", c.testResponses.recponseEvent.Status)
	}
	return nil
}

func (c *calendarTest) iReceiveEventWithIDAndUserIDTitleDescription(userID int, title, description string) error {
	if uint64(userID) != c.testResponses.recponseEvent.Item.UserId {
		return errors.Errorf("Expected userID: %v, have received userID: %v", userID, c.testResponses.recponseEvent.Item.UserId)
	}

	if title != c.testResponses.recponseEvent.Item.Title {
		return errors.Errorf("Expected userID: %v, have received userID: %v", title, c.testResponses.recponseEvent.Item.Title)
	}

	if description != c.testResponses.recponseEvent.Item.Description {
		return errors.Errorf("Expected userID: %v, have received userID: %v", description, c.testResponses.recponseEvent.Item.Description)
	}

	return nil
}

func (c *calendarTest) iCreateEventWithUserIDTitle(userID int, title string) error {
	event := event.Event{
		UserID:       uint(userID),
		Title:        title,
		Time:         time.Now().AddDate(1, 0, 1),
	}
	e, err := controller.EventToEventProto(event)
	if err != nil {
		return err
	}

	_, err = c.GRPCClient.EventCreate(c.ctx, e)
	if err != nil {
		c.testResponses.error = err
	}

	return nil
}

func (c *calendarTest) iReceiveStatusIsNotOK() error {
	if c.testResponses.error == nil {
		return errors.Errorf("Expected error, have received error = nil")
	}
	return nil
}

func (c *calendarTest) iCreateEventOnANextDay() error {
	event := event.Event{
		UserID:       uint(userID),
		Title:        "EventOnANextDay",
		Time:         time.Now().AddDate(0, 0, 1),
	}
	e, err := controller.EventToEventProto(event)
	if err != nil {
		return err
	}

	response, err := c.GRPCClient.EventCreate(c.ctx, e)
	if err != nil {
		return err
	}
	c.testResponses.recponseEvent = response

	return nil
}

func (c *calendarTest) iSendRequestForAListOfEventsOnADay() error {
	response, err := c.GRPCClient.EventListOnDay(c.ctx, ptypes.TimestampNow())
	if err != nil {
		return err
	}
	c.testResponses.recponseEvents = response

	return nil
}

func (c *calendarTest) iReceiveListOfEventsWithLengthEqualToOne() error {
	length := 1
	if len(c.testResponses.recponseEvents.List) != length {
		return errors.Errorf("Expected length %v, have received lehgth: %v", length, len(c.testResponses.recponseEvents.List))
	}
	return nil
}

func (c *calendarTest) iCreateEventOnASecondDay() error {
	event := event.Event{
		UserID:       uint(userID),
		Title:        "EventOnASecondDay",
		Time:         time.Now().AddDate(0, 0, 2),
	}
	e, err := controller.EventToEventProto(event)
	if err != nil {
		return err
	}

	response, err := c.GRPCClient.EventCreate(c.ctx, e)
	if err != nil {
		return err
	}
	c.testResponses.recponseEvent = response

	return nil
}

func (c *calendarTest) iSendRequestForAListOfEventsOnAWeek() error {
	response, err := c.GRPCClient.EventListOnWeek(c.ctx, ptypes.TimestampNow())
	if err != nil {
		return err
	}

	c.testResponses.recponseEvents = response

	return nil
}

func (c *calendarTest) iReceiveListOfEventsWithLengthEqualToTwo() error {
	length := 2
	if len(c.testResponses.recponseEvents.List) != length {
		return errors.Errorf("Expected length %v, have received lehgth: %v", length, len(c.testResponses.recponseEvents.List))
	}
	return nil
}

func (c *calendarTest) iCreateEventOnASecondWeek() error {
	event := event.Event{
		UserID:       uint(userID),
		Title:        "EventOnASecondWeek",
		Time:         time.Now().AddDate(0, 0, 9),
	}
	e, err := controller.EventToEventProto(event)
	if err != nil {
		return err
	}

	response, err := c.GRPCClient.EventCreate(c.ctx, e)
	if err != nil {
		return err
	}
	c.testResponses.recponseEvent = response

	return nil
}

func (c *calendarTest) iSendRequestForAListOfEventsOnAMonth() error {
	response, err := c.GRPCClient.EventListOnMonth(c.ctx, ptypes.TimestampNow())
	if err != nil {
		return err
	}

	c.testResponses.recponseEvents = response

	return nil
}

func (c *calendarTest) iReceiveListOfEventsWithLengthEqualToThree() error {
	length := 3
	if len(c.testResponses.recponseEvents.List) != length {
		return errors.Errorf("Expected length %v, have received lehgth: %v", length, len(c.testResponses.recponseEvents.List))
	}
	return nil
}

func (c *calendarTest) iCreateEventOnANextDayWithDurationInDayAndTitle(title string) error {
	event := event.Event{
		UserID:     uint(userID),
		Title:      title,
		Duration:	25 * time.Hour,
		Time:       time.Now().AddDate(0, 0, 1),
	}
	e, err := controller.EventToEventProto(event)
	if err != nil {
		return err
	}

	response, err := c.GRPCClient.EventCreate(c.ctx, e)
	if err != nil {
		return err
	}
	c.testResponses.recponseEvent = response

	return nil
}

func (c *calendarTest) iReceiveNotificationWithIDAndTitle(title string) error {
	n := <- c.queueCh

	if n.Title != title {
		return errors.Errorf("Expected title of notice: %q, have received title: %q", title, n.Title)
	}

	return nil
}

func FeatureContext(s *godog.Suite) {
	c := newCalendarTest(context.Background())

	s.BeforeSuite(c.start)

	s.Step(`^I create event with UserID=(\d+), Title="([^"]*)", Description="([^"]*)"$`, c.iCreateEventWithUserIDTitleDescription)
	s.Step(`^I receive status is OK$`, c.iReceiveStatusIsOK)
	s.Step(`^I receive event with ID and UserID=(\d+), Title="([^"]*)", Description="([^"]*)"$`, c.iReceiveEventWithIDAndUserIDTitleDescription)
	s.Step(`^I create event with UserID=(\d+), Title="([^"]*)"$`, c.iCreateEventWithUserIDTitle)
	s.Step(`^I receive status is not OK$`, c.iReceiveStatusIsNotOK)
	s.Step(`^I create event on a next day$`, c.iCreateEventOnANextDay)
	s.Step(`^I send request for a list of events on a day$`, c.iSendRequestForAListOfEventsOnADay)
	s.Step(`^I receive list of events with length  equal to one$`, c.iReceiveListOfEventsWithLengthEqualToOne)
	s.Step(`^I create event on a second day$`, c.iCreateEventOnASecondDay)
	s.Step(`^I send request for a list of events on a week$`, c.iSendRequestForAListOfEventsOnAWeek)
	s.Step(`^I receive list of events with length  equal to two$`, c.iReceiveListOfEventsWithLengthEqualToTwo)
	s.Step(`^I create event on a second week$`, c.iCreateEventOnASecondWeek)
	s.Step(`^I send request for a list of events on a month$`, c.iSendRequestForAListOfEventsOnAMonth)
	s.Step(`^I receive list of events with length  equal to three$`, c.iReceiveListOfEventsWithLengthEqualToThree)
	s.Step(`^I create event on a next day with duration in day and Title="([^"]*)",$`, c.iCreateEventOnANextDayWithDurationInDayAndTitle)
	s.Step(`^I receive notification with ID and Title="([^"]*)"$`, c.iReceiveNotificationWithIDAndTitle)

	s.AfterSuite(c.stop)

}



