package main

import (
	"context"
	"flag"
	"github.com/Kalinin-Andrey/otus-go/hw10/pkg/rw"
	"log"
	"net"
	"os"
	"strings"
	"sync"
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

	sincStop := make(chan struct{}, 1)
	defer close(sincStop)
	wg := &sync.WaitGroup{}
	wg.Add(2)
	ctx := rw.StopSynchronizer(context.Background(), wg, sincStop)

	conn, err := net.DialTimeout("tcp", pars[0] + ":" + pars[1], timeout)
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer conn.Close()
	go rw.ReadRoutine(ctx, wg, conn, os.Stdout, func(){
		sincStop <- struct{}{}
		conn.Write([]byte("\n"))
	})
	go rw.WriteRoutine(ctx, wg, conn, os.Stdin, func(){
		sincStop <- struct{}{}
	})
	wg.Wait()
	log.Println("Client completed work")
}
