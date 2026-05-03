package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	desc "github.com/malfoit/SimpleProject/pkg/user/v1"
	"google.golang.org/grpc"
)

func main() {
	c := newContainer()

	addr := fmt.Sprintf(":%s", c.config.GRPC.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	desc.RegisterUserV1Server(grpcServer, c.userHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("gRPC server listening on %s", addr)
		if err = grpcServer.Serve(listener); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down gRPC server...")
	grpcServer.GracefulStop()
	log.Println("server stopped")
}
