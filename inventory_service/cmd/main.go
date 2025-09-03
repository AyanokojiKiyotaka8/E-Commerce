package main

import (
	"context"
	"log"

	"github.com/AyanokojiKiyotaka8/E-Commerce/inventory_service/internal/kafka"
	"github.com/AyanokojiKiyotaka8/E-Commerce/inventory_service/internal/service"
	"github.com/AyanokojiKiyotaka8/E-Commerce/inventory_service/internal/store"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017/"))
	if err != nil {
		log.Fatal(err)
	}

	inventoryStore := store.NewMongoInventoryStore(client)
	inventoryService := service.NewInventoryService(inventoryStore)
	kafkaConsumer, err := kafka.NewKafkaConsumer(inventoryService)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaConsumer.Stop()
	kafkaConsumer.Start()
}
