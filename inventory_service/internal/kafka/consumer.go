package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/AyanokojiKiyotaka8/E-Commerce/inventory_service/internal/service"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/sirupsen/logrus"
)

var producerTopics = []string{"product-created", "product-updated-old", "product-updated-new", "product-deleted"}

type Consumer interface {
	Consume()
	Start()
	Stop()
}

type KafkaConsumer struct {
	consumer  *kafka.Consumer
	isRunning bool
	svc       service.InventoryServicer
}

func NewKafkaConsumer(svc service.InventoryServicer) (*KafkaConsumer, error) {
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("kafka consumer creation problem")
	}

	if err := consumer.SubscribeTopics(producerTopics, nil); err != nil {
		return nil, fmt.Errorf("failed to subscribe to topics %w", err)
	}
	return &KafkaConsumer{
		consumer: consumer,
		svc:      svc,
	}, nil
}

func (c *KafkaConsumer) Start() {
	c.isRunning = true
	c.Consume()
}

func (c *KafkaConsumer) Consume() {
	for c.isRunning {
		message, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error: %s", err)
			continue
		}

		topic := message.TopicPartition.Topic
		var stockKey string
		if err := json.Unmarshal(message.Value, &stockKey); err != nil {
			logrus.Errorf("JSON serialization error: %s", err)
			continue
		}

		if err := c.svc.GetConsumeData(*topic, stockKey); err != nil {
			logrus.WithFields(logrus.Fields{
				"topic":    *topic,
				"stockKey": stockKey,
				"error":    err,
			}).Error("Failed to handle product data")
			continue
		}
	}
}

func (c *KafkaConsumer) Stop() {
	c.isRunning = false
	c.consumer.Close()
}
