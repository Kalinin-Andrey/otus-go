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


// ReadRoutine func is read from conn and write to out
func ReadRoutine(ctx context.Context, wg *sync.WaitGroup, conn net.Conn, out io.Writer, stopWriteRoutine func()) {
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
				stopWriteRoutine()
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
	log.Printf("ReadRoutine has finished\n")
}

// WriteRoutine ir read from in and write in conn
func WriteRoutine(ctx context.Context, wg *sync.WaitGroup, conn net.Conn, in io.Reader, stopReadRoutine func()) {
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
				stopReadRoutine()
				break OUTER
			}
			log.Printf("scan\n")
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
	log.Printf("WriteRoutine has finished\n")
}

// StopSynchronizer synchronize a stoppage of script
func StopSynchronizer(ctx context.Context, wg *sync.WaitGroup, sincStop chan struct{}, sig ...os.Signal) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	ctx, finish := context.WithCancel(ctx)
	go func() {
		defer wg.Done()
		defer close(c)
	OUTER:
		for {
			select {
			case s := <- c:
			if s != nil {
				log.Printf("Got siognal: %v\n", s)
				break OUTER
			}
			case <- sincStop:
				log.Printf("Got sincStop siognal\n")
				break OUTER
			}
		}
		finish()
		log.Println("StopSynchronizer has finished")
	}()
	return ctx
}