package rabbitmq

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	"github.com/cenkalti/backoff/v3"
	"github.com/streadway/amqp"

	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/config"
	"github.com/Kalinin-Andrey/otus-go/calendar/pkg/log"
)

type QueueClient interface {
	Publish(body []byte) error
	Handle(fn func(<-chan amqp.Delivery), threads int) error
	Close() error
}

// Client ...
type Client struct {
	ctx				context.Context
	clientType			uint
	logger			log.ILogger
	conn            *amqp.Connection
	channel         *amqp.Channel
	reconnected     chan error
	consumerTag     string
	uri             string
	exchangeName    string
	exchangeType    string
	queue           string
	bindingKey      string
	deliveryChannel <-chan amqp.Delivery
}

const (
	TypePublisher	= 0
	TypeConsumer	= 1
)

var _ QueueClient = (*Client)(nil)

var defaultExchangeType = amqp.ExchangeDirect

func NewClient(ctx context.Context, logger log.ILogger, conf config.RabbitMQ, clientType uint) (*Client, error) {
	exchangeType	:= conf.ExchangeType
	if conf.ExchangeType == "" {
		exchangeType	= defaultExchangeType
	}


	c := &Client{
		ctx:			ctx,
		clientType:		clientType,
		logger:			logger.With(ctx),
		consumerTag:	conf.ConsumerTag,
		uri:			conf.Uri,
		exchangeName:	conf.ExchangeName,
		exchangeType:	exchangeType,
		queue:			conf.Queue,
		bindingKey:		conf.BindingKey,
		reconnected:	make(chan error),
	}
	return c, c.connect()
}

func (c *Client) getDeliveryChannel() <-chan amqp.Delivery {
	return c.deliveryChannel
}

func (c *Client) getAmqpChannel() *amqp.Channel {
	return c.channel
}

func (c *Client) reConnect() error {
	be := backoff.NewExponentialBackOff()
	be.MaxElapsedTime = time.Minute
	be.InitialInterval = 1 * time.Second
	be.Multiplier = 2
	be.MaxInterval = 15 * time.Second

	b := backoff.WithContext(be, context.Background())
	for {
		d := b.NextBackOff()
		if d == backoff.Stop {
			return fmt.Errorf("stop reconnecting")
		}

		select {
		case <-c.ctx.Done():
			break
		case <-time.After(d):
			err := c.connect()
			if err != nil {
				fmt.Printf("Couldn't connect: %+v", err)
				continue
			}

			return nil
		}
	}
}

func (c *Client) connect() error {
	var err error

	c.conn, err = amqp.Dial(c.uri)
	if err != nil {
		return errors.Errorf("Dial err: %v", err)
	}

	go func() {
		select {
		case <- c.ctx.Done():
			c.Close()
			return
		case err := <- c.conn.NotifyClose(make(chan *amqp.Error)):
			c.logger.Infof("Reconnecting Error: %s", err)
		}
		// Понимаем, что канал сообщений закрыт, надо пересоздать соединение.
		err = c.reConnect()
		if err != nil {
			c.logger.Errorf("Reconnecting Error: %s", err)
			return
		}

		select {
		case c.reconnected <- errors.New("the channel was reconnected"):
		default:
		}
		return
	}()

	c.channel, err = c.conn.Channel()
	if err != nil {
		return errors.Errorf("Channel err: %v", err)
	}

	if err = c.channel.ExchangeDeclare(
		c.exchangeName,
		c.exchangeType,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return errors.Errorf("ExchangeDeclare err: %v", err)
	}

	queue, err := c.channel.QueueDeclare(
		c.queue,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		return errors.Errorf("QueueDeclare err: %v", err)
	}

	// Число сообщений, которые можно подтвердить за раз.
	err = c.channel.Qos(50, 0, false)
	if err != nil {
		return errors.Errorf("error setting qos err: %v", err)
	}

	// Создаём биндинг (правило маршрутизации).
	if err = c.channel.QueueBind(
		queue.Name,
		c.bindingKey,
		c.exchangeName,
		false,
		nil,
	); err != nil {
		return errors.Errorf("QueueBind err: %v", err)
	}

	if c.clientType == TypeConsumer {
		ch, err := c.channel.Consume(
			queue.Name,
			c.consumerTag,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return errors.Errorf("Consume err: %v", err)
		}
		c.deliveryChannel = ch
	}
	return nil
}

func (c *Client) Handle(fn func(<-chan amqp.Delivery), threads int) error {

	go func() {
		for {
			for i := 0; i < threads; i++ {
				go fn(c.getDeliveryChannel())
			}

			select {
			case <-c.ctx.Done():
				return
			case err := <-c.reconnected:
				c.logger.Infof("[rabbitmq] reconnected: ", err)
			}
		}
	}()
	return nil
}


func (c *Client) Close() error {
	if c.conn == nil {
		return nil
	}
	return c.conn.Close()
}


func (c *Client) Publish(body []byte) error {

	err := c.getAmqpChannel().Publish(
		c.exchangeName,     // exchange
		c.bindingKey, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
	if err != nil {
		return errors.Wrapf(err, "Failed to publish a message")
	}
	c.logger.Debugf(" [x] Sent %s", body)
	return nil
}


