package server

import (
	"context"
	"time"

	"53HW/db"
	pb "53HW/proto"
)

type ChatServer struct {
	pb.UnimplementedChatServiceServer
}

func (s *ChatServer) SendMessage(ctx context.Context, req *pb.ChatMessage) (*pb.Empty, error) {
	_, err := db.DB.Exec("INSERT INTO messages (username, message, timestamp) VALUES ($1, $2, $3)",
		req.Username, req.Message, time.Now())
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (s *ChatServer) StreamMessages(req *pb.Empty, stream pb.ChatService_StreamMessagesServer) error {
	rows, err := db.DB.Query("SELECT username, message, timestamp FROM messages ORDER BY timestamp ASC")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var msg pb.ChatMessage
		if err := rows.Scan(&msg.Username, &msg.Message, &msg.Timestamp); err != nil {
			return err
		}
		if err := stream.Send(&msg); err != nil {
			return err
		}
	}
	return nil
}
