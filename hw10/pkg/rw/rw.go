package rw

import (
	"context"
	"io"
	"bufio"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)



func ReadRoutine(ctx context.Context, wg *sync.WaitGroup, conn net.Conn, out io.Writer) {
	defer wg.Done()
	scanner := bufio.NewScanner(conn)
	//writer := bufio.NewWriter(out)
OUTER:
	for {
		select {
		case <-ctx.Done():
			log.Printf("ReadRoutine: ctx.Done()\n")
			break OUTER
		default:
			if !scanner.Scan() {
				log.Printf("CANNOT SCAN\n")
				break OUTER
			}
			s := scanner.Text()
			log.Printf("ReadRoutine read: %v\n", s)
			io.WriteString(out, s + "\n")
			//writer.WriteString(s)

			if err := scanner.Err(); err != nil {
				log.Printf("ReadRoutine: error happend: %v\n", err)
			}
		}
	}
	log.Printf("Finished ReadRoutine\n")
}

func WriteRoutine(ctx context.Context, wg *sync.WaitGroup, conn net.Conn, in io.Reader) {
	defer wg.Done()
	scanner := bufio.NewScanner(in)
	//writer := bufio.NewWriter(conn)
OUTER:
	for {
		select {
		case <-ctx.Done():
			log.Printf("ReadRoutine: ctx.Done()\n")
			break OUTER
		default:
			if !scanner.Scan() {
				log.Printf("CANNOT SCAN\n")
				break OUTER
			}
			s := scanner.Text()

			conn.Write([]byte(s + "\n"))
			log.Printf("WriteRoutine write: %v\n", s)
			//io.WriteString(conn, s)
			//writer.WriteString(s)

			if err := scanner.Err(); err != nil {
				log.Printf("WriteRoutine: error happend: %v\n", err)
			}
		}

	}
	log.Printf("Finished WriteRoutine\n")
}

func ContextWithCancelBySignal(ctx context.Context, sig ...os.Signal) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	ctx, finish := context.WithCancel(ctx)
	go func() {
		defer close(c)
		for {
			s := <-c
			if s != nil {
				log.Printf("Got siognal: %v", s)
				finish()
				break
			}
		}
	}()
	return ctx
}