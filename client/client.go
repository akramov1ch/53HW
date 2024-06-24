package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"53HW/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:8001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := proto.NewChatServiceClient(conn)

	sendMessage(client, "user1", "Hello, world!")

	streamMessages(client)
}

func sendMessage(client proto.ChatServiceClient, username, message string) {
	req := &proto.ChatMessage{
		Username:  username,
		Message:   message,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	_, err := client.SendMessage(context.Background(), req)
	if err != nil {
		log.Fatalf("Could not send message: %v", err)
	}

	fmt.Println("Message sent successfully!")
}

func streamMessages(client proto.ChatServiceClient) {
	stream, err := client.StreamMessages(context.Background(), &proto.Empty{})
	if err != nil {
		log.Fatalf("Could not stream messages: %v", err)
	}

	fmt.Println("Streaming messages:")
	for {
		msg, err := stream.Recv()
		if err != nil {
			log.Fatalf("Error receiving message: %v", err)
		}
		fmt.Printf("[%s] %s: %s\n", msg.Timestamp, msg.Username, msg.Message)
	}
}
