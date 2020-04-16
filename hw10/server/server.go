package main

import (
	"context"
	"fmt"
	"github.com/Kalinin-Andrey/otus-go/hw10/pkg/rw"
	"log"
	"net"
	"os"
	"sync"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte(fmt.Sprintf("Welcome to %s, friend from %s\n", conn.LocalAddr(), conn.RemoteAddr())))

	sincStop := make(chan struct{}, 1)
	defer close(sincStop)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx := rw.StopSynchronizer(context.Background(), wg, sincStop)
	go rw.ReadRoutine(ctx, wg, conn, os.Stdout, func(){
		sincStop <- struct{}{}
		conn.Write([]byte("\n"))
	})
	go rw.WriteRoutine(ctx, wg, conn, os.Stdin, func(){
		sincStop <- struct{}{}
	})
	wg.Wait()

	log.Printf("Closing connection with %s", conn.RemoteAddr())

}

func main() {
	l, err := net.Listen("tcp", "0.0.0.0:3301")
	if err != nil {
		log.Fatalf("Cannot listen: %v", err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("Cannot accept: %v", err)
		}

		go handleConnection(conn)
	}
}
