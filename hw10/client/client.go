package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"time"
)

const (
	defaultTimeout = "10s"
)

var sTimeout string


func init() {
	flag.StringVar(&sTimeout, "timeout", defaultTimeout, "timeout")
}



func main() {
	flag.Parse()
	pars := make([]string, 0, 2)

	for i := 1; i < len(os.Args); i++ {
		if !strings.HasPrefix(os.Args[i], "--timeout") {
			pars = append(pars, os.Args[i])
		}
	}
	timeout, err := time.ParseDuration(sTimeout)
	if err != nil {
		log.Fatal(err)
	}
	if len(pars) < 2 {
		log.Fatal("requared params: host port")
	}

	c := NewClient(context.Background(),"tcp", pars[0] + ":" + pars[1], timeout, os.Stdin, os.Stdout)

	c.Run()
}

type client struct {
	network		string
	address		string
	timeout		time.Duration
	in			io.Reader
	out			io.Writer
	conn		net.Conn
	sincStop	chan struct{}
	ctx			context.Context
	finish		context.CancelFunc
}

func NewClient(ctx context.Context, network string, address string, timeout time.Duration, in io.Reader, out io.Writer) *client {
	return &client{
		network:		network,
		address:		address,
		timeout:		timeout,
		in:				in,
		out:			out,
		sincStop:		make(chan struct{}, 1),
		ctx:			ctx,
	}
}

// Run a client
func (c *client) Run() {
	defer close(c.sincStop)
	c.StopSynchronizer()

	conn, err := net.DialTimeout(c.network, c.address, c.timeout)
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	c.conn = conn
	defer c.conn.Close()

	go c.readRoutine()
	go c.writeRoutine()
	<- c.ctx.Done()
	log.Println("Client completed work")
}

// Stop a client
func (c *client) Stop() {
	c.sincStop <- struct{}{}
}

// readRoutine func reads from conn and writes to out
func (c *client) readRoutine() {
	scanner := bufio.NewScanner(c.conn)
OUTER:
	for {
		select {
		case <-c.ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				c.Stop()
				break OUTER
			}
			s := scanner.Text()
			io.WriteString(c.out, s + "\n")

			if err := scanner.Err(); err != nil {
				log.Printf("ReadRoutine: error happend: %v\n", err)
			}
		}
	}
	log.Printf("ReadRoutine has finished\n")
}

// writeRoutine func reads from in and writes in conn
func (c *client) writeRoutine() {
	scanner := bufio.NewScanner(c.in)
OUTER:
	for {
		select {
		case <-c.ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				c.Stop()
				break OUTER
			}
			s := scanner.Text()
			_, err := c.conn.Write([]byte(s + "\n"))

			if err != nil {
				log.Printf("GRPC write error: %v\n", err)
			}

			if err = scanner.Err(); err != nil {
				log.Printf("Scanner error: %v\n", err)
			}
		}

	}
	log.Printf("WriteRoutine has finished\n")
}

// StopSynchronizer synchronize a stoppage of script
func (c *client) StopSynchronizer() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch)
	c.ctx, c.finish = context.WithCancel(c.ctx)
	go func() {
		defer func() {
			signal.Stop(ch)
			//close(c) // ?
		}()
	OUTER:
		for {
			select {
			case s := <- ch:
				if s != nil {
					break OUTER
				}
			case <- c.sincStop:
				break OUTER
			}
		}
		c.finish()
		log.Println("StopSynchronizer has finished")
	}()
}