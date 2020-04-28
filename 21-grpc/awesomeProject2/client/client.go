package main

import (
	chat "awesomeProject2/chatpb"
	"bufio"
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"

	"os"
	"time"
)

func writeRoutine(end chan interface{}, ctx context.Context, conn chat.ChatExampleClient) {
	scanner := bufio.NewScanner(os.Stdin)
OUTER:
	for {
		select {
		case <-ctx.Done():
			break OUTER
		default:
			if !scanner.Scan() {
				break OUTER
			}
			str := scanner.Text()

			if str == "end" {
				break OUTER
			}
			log.Printf("To server %v\n", str)

			msg, err := conn.SendMessage(context.Background(), &chat.ChatMessage{
				Text:    str,
				Created: ptypes.TimestampNow(),
			})

			if err != nil {
				errMsg := status.Convert(err)
				fmt.Printf("err %s %s", errMsg.Code(), errMsg.Message())
			}

			if msg != nil {
				created, _ := ptypes.Timestamp(msg.Created)
				created = created.Local()
				fmt.Printf("%s created %s", msg.Text, created)
			}

		}
	}

	log.Printf("Finished writeRoutine")
	close(end)
}

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Minute)

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := chat.NewChatExampleClient(cc)
	end := make(chan interface{})
	go writeRoutine(end, ctx, c)

	<-end
}
