package main

import (
	"log"
	"net"

	"github.com/AyanokojiKiyotaka8/E-Commerce/types"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	productService := NewProductService()
	grpcServer := grpc.NewServer()
	types.RegisterProductServer(grpcServer, productService)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
