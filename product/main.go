package main

import (
	"context"
	"log"
	"net"

	"github.com/AyanokojiKiyotaka8/E-Commerce/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	productStore := NewMongoProductStore(client)

	productService := NewProductService(productStore)
	grpcServer := grpc.NewServer()
	types.RegisterProductServiceServer(grpcServer, productService)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
