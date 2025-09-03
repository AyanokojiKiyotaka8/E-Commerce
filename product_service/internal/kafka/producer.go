package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Producer interface {
	Produce(string, string) error
	Stop()
}

type KafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer() (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost"})
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error == nil {
					log.Printf("Delivered to %v\n", ev.TopicPartition)
				} else {
					log.Printf("Delivery failed: %v\n", ev.TopicPartition.Error)
				}
			}
		}
	}()

	return &KafkaProducer{producer: p}, nil
}

func (p *KafkaProducer) Produce(topic string, event string) error {
	ev, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	err = p.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: ev,
	}, nil)
	if err != nil {
		return fmt.Errorf("failed to produce message: %w", err)
	}
	return nil
}

func (p *KafkaProducer) Stop() {
	p.producer.Flush(15_000)
	p.producer.Close()
}
