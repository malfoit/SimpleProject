package main

import (
	"context"
	"fmt"
	"log"
	"net"

	desc "github.com/malfoit/SimpleProject/pkg/user/v1"
	"google.golang.org/grpc"
)

func main() {
	container := NewContainer()

	ctx := context.Background()

	synced, err := container.SyncService.SyncUsers(ctx)
	if err != nil {
		log.Printf("Sync failed: %v", err)
	} else {
		log.Printf("Synced %d users", synced)
	}

	grpcServer := grpc.NewServer()
	desc.RegisterUserV1Server(grpcServer, container.UserHandler)

	addr := fmt.Sprintf(":%s", container.Config.GRPC.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("Server started on %s", addr)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
