package main

import (
	chat "awesomeProject2/chatpb"
	"bufio"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte(fmt.Sprintf("Welcome to %s, friend from %s\n", conn.LocalAddr(), conn.RemoteAddr())))

	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := scanner.Text()
		log.Printf("RECEIVED: %s", text)
		if text == "quit" || text == "exit" {
			break
		}

		conn.Write([]byte(fmt.Sprintf("I have received '%s'\n", text)))
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error happend on connection with %s: %v", conn.RemoteAddr(), err)
	}

	log.Printf("Closing connection with %s", conn.RemoteAddr())

}

type ServerExample struct {
}

var i int64

func (s ServerExample) SendMessage(context context.Context, msg *chat.ChatMessage) (*chat.ChatMessage, error) {
	defer func() { i++ }()
	if msg.Text == "" {
		return nil, status.Error(codes.InvalidArgument,"no empty string")
	}
	return &chat.ChatMessage{Id: i, Text: "Pong:" + msg.Text,Created: ptypes.TimestampNow(),}, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen %v", err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	chat.RegisterChatExampleServer(grpcServer, ServerExample{})
	_ = grpcServer.Serve(lis)
}
