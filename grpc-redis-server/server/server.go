package main

import (
	"fmt"
	"log"
	"net"

	pb "grpc-redis-server/protos/todo/protos/todo"

	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc"
)

func main() {
	funcLogs := hclog.Default()

	lis, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	todoService := NewTodoServiceServer("localhost:6379", funcLogs)

	pb.RegisterTodoServiceServer(grpcServer, todoService)

	fmt.Println("Running Server on Port 8888...")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
