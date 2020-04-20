package rw

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
)


// ReadRoutine func is read from conn and write to out
func ReadRoutine(ctx context.Context, conn net.Conn, out io.Writer, stopWriteRoutine func()) {
	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		s := scanner.Text()
		log.Printf("ReadRoutine read: %v\n", s)
		io.WriteString(out, s + "\n")

		if err := scanner.Err(); err != nil {
			log.Printf("ReadRoutine: error happend: %v\n", err)
		}
	}
	//stopWriteRoutine()
	log.Printf("ReadRoutine has finished\n")
}

// WriteRoutine ir read from in and write in conn
func WriteRoutine(ctx context.Context, conn net.Conn, in io.Reader, stopReadRoutine func()) {
	scanner := bufio.NewScanner(in)

	for scanner.Scan() {
		log.Printf("scan\n")
		s := scanner.Text()
		conn.Write([]byte(s + "\n"))

		log.Printf("WriteRoutine write: %v\n", s)

		if err := scanner.Err(); err != nil {
			log.Printf("WriteRoutine: error happend: %v\n", err)
		}

	}
	//stopReadRoutine()
	log.Printf("WriteRoutine has finished\n")
}

// StopSynchronizer synchronize a stoppage of script
func StopSynchronizer(ctx context.Context, sincStop chan struct{}, sig ...os.Signal) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, sig...)
	ctx, finish := context.WithCancel(ctx)
	go func() {
		defer func() {
			signal.Stop(c)
			close(c) // ?
		}()
	OUTER:
		for {
			select {
			case s := <- c:
			if s != nil {
				//log.Printf("Got siognal: %v\n", s)
				break OUTER
			}
			case <- sincStop:
				//log.Printf("Got sincStop siognal\n")
				break OUTER
			}
		}
		finish()
		log.Println("StopSynchronizer has finished")
	}()
	return ctx
}