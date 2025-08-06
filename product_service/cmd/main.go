package main

import (
	"context"
	"log"
	"net"

	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/handler"
	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/kafka"
	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/middleware"
	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/service"
	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/internal/store"
	"github.com/AyanokojiKiyotaka8/E-Commerce/product_service/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	productStore := store.NewMongoProductStore(client)
	kafkaProducer, err := kafka.NewKafkaProducer()
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaProducer.Stop()

	var productService service.ProductServicer
	productService = service.NewProductService(productStore, kafkaProducer)
	productService = middleware.NewLogMiddleware(productService)
	productHandler := handler.NewProductHandler(productService)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterProductServiceServer(grpcServer, productHandler)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
