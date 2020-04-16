package main

import (
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "0.0.0.0:3303")
	if err != nil {
		log.Fatalf("Can't resolve addr: %v", err)
	}

	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("Can't listen: %v", err)
	}
	defer l.Close()

	for {
		msg := make([]byte, 1024)
		length, fromAddr, err := l.ReadFromUDP(msg)
		if err != nil {
			log.Fatalf("Error happened: %v", err)
		}

		log.Printf("Message from %s with length %d: %s\n", fromAddr.String(), length, string(msg))
	}
}
