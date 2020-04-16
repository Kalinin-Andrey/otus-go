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


	ctx := rw.ContextWithCancelBySignal(context.Background())

	conn, err := net.DialTimeout("tcp", pars[0] + ":" + pars[1], timeout)
	if err != nil {
		log.Fatalf("Cannot connect: %v", err)
	}
	defer conn.Close()
	wg := &sync.WaitGroup{}
	wg.Add(2)
	go rw.ReadRoutine(ctx, wg, conn, os.Stdout)
	go rw.WriteRoutine(ctx, wg, conn, os.Stdin)
	wg.Wait()
	log.Println("Client completed work")
}
