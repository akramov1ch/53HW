package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"53HW/config"
	"53HW/db"
	pb "53HW/proto"
	"53HW/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	err = db.InitDB(cfg)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	err = db.CreateMessagesTable()
	if err != nil {
		log.Fatalf("Could not create messages table: %v", err)
	}

	lis, err := net.Listen("tcp", "8001")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServiceServer(s, &server.ChatServer{})

	log.Println("Server is running on port 50051...")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
