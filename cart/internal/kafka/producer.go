package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	flushTimeout = 5000
)

type (
	CartEventProducer interface {
		ProduceCartItemAdded(ctx context.Context, payload CartItemAddedPayload)
		ProduceCartItemFailed(ctx context.Context, payload CartItemFailedPayload)
		produce(ctx context.Context, message []byte, key string, partition int32)
		Close()
	}
)

// Cart item event models.
type (
	EventModel struct {
		Type      string      `json:"type"`
		Service   string      `json:"service"`
		Timestamp time.Time   `json:"timestamp"`
		Payload   interface{} `json:"payload"`
	}

	CartItemAddedPayload struct {
		CartID string `json:"cartId"`
		SKU    uint32 `json:"sku"`
		Count  uint16 `json:"count"`
		Status string `json:"status"`
	}

	CartItemFailedPayload struct {
		CartID string `json:"cartId"`
		SKU    uint32 `json:"sku"`
		Count  uint16 `json:"count"`
		Status string `json:"status"`
		Reason string `json:"reason"`
	}
)

var _ CartEventProducer = (*cartEventProducer)(nil)

type cartEventProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewCartServiceProducer(address []string) (*cartEventProducer, error) {
	conf := &kafka.ConfigMap{
		"bootstrap.servers": strings.Join(address, ","),
	}

	prod, err := kafka.NewProducer(conf)
	if err != nil {
		return nil, fmt.Errorf("error creating cart service kafka producer: %w", err)
	}

	// i started goroutine for delivery report here.
	go func() {
		for e := range prod.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("[cartService]delivery failed for topic %s [partition %d]: %v\n", *ev.TopicPartition.Topic,
						ev.TopicPartition.Partition, ev.TopicPartition.Error,
					)
				} else {
					log.Printf("[cartService]delivered message to topic %s [partition %d]: %v\n", *ev.TopicPartition.Topic,
						ev.TopicPartition.Partition, ev.TopicPartition.Error,
					)
				}
			case kafka.Error:
				log.Printf("cart service kafka producer error: %v\n", ev)
			}
		}
	}()

	return &cartEventProducer{
		producer: prod,
		topic:    "metrics",
	}, nil
}

func (cp *cartEventProducer) ProduceCartItemAdded(ctx context.Context, payload CartItemAddedPayload) {
	event := EventModel{
		Type:      "cart_item_added",
		Service:   "cart",
		Timestamp: time.Now(),
		Payload:   payload,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal cart_item_added event: %v\n", err.Error())
	}

	cp.produce(ctx, eventBytes, "some_key", 0)
}

func (cp *cartEventProducer) ProduceCartItemFailed(ctx context.Context, payload CartItemFailedPayload) {
	event := EventModel{
		Type:      "cart_item_failed",
		Service:   "cart",
		Timestamp: time.Now(),
		Payload:   payload,
	}

	eventBytes, err := json.Marshal(event)
	if err != nil {
		log.Printf("failed to marshal cart_item_failed event: %v\n", err.Error())
	}

	cp.produce(ctx, eventBytes, "another_key", 0)
}

func (cp *cartEventProducer) produce(ctx context.Context, message []byte, key string, partition int32) {
	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &cp.topic,
			Partition: partition,
		},
		Value:     message,
		Key:       []byte(key),
		Timestamp: time.Now(),
	}

	kafkaChan := make(chan kafka.Event, 1)
	if err := cp.producer.Produce(kafkaMessage, kafkaChan); err != nil {
		log.Printf("failed to produce message with key %s: %v\n", key, err)
	}

	go func() {
		select {
		case e := <-kafkaChan:
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Printf("[cartService_kafkaProducer]: delivery failed for key %s: %v\n", key,
						ev.TopicPartition.Error,
					)
				} else {
					log.Printf("[cartService_kafkaProducer]: delivery message with key %s to topic %s [partition %d] at offset %v\n",
						key, *ev.TopicPartition.Topic, ev.TopicPartition.Partition, ev.TopicPartition.Offset,
					)
				}
			case kafka.Error:
				log.Printf("[cartService_kafkaProducer]: error for key %s: %v\n", key, ev)
			}
		case <-ctx.Done():
			log.Printf("context cancelled for message with key %s: %v\n", key, ctx.Err())
		}
	}()
}

func (cp *cartEventProducer) Close() {
	cp.producer.Flush(flushTimeout)
	cp.producer.Close()
}
