package kafka

import (
	"context"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	sessionTimeoutMs = 7000
	timeOut          = 2
)

type Handler interface {
	HandleMessage(message []byte, topic kafka.Offset) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
	stop     bool
}

func NewConsumer(handler Handler, address string, topic, consumerGroup string) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":        address,
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeoutMs,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create new metrics consumer: %w", err)
	}

	if err := c.Subscribe(topic, nil); err != nil {
		return nil, fmt.Errorf("[metrics-consumer]: c.Subscribe: %w", err)
	}

	return &Consumer{
		consumer: c,
		handler:  handler,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("context cancelled, stopping consumer...")
			return
		default:
			kafkaMsg, err := c.consumer.ReadMessage(timeOut)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() != kafka.ErrTimedOut {
					log.Printf("c.consumer.ReadMessage: error: %v\n", err.Error())
				}

				continue
			}

			if kafkaMsg == nil {
				continue
			}

			if err := c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset); err != nil {
				log.Printf("c.handler.HandleMessage: %v\n", err.Error())
				continue
			}

			if _, err := c.consumer.StoreMessage(kafkaMsg); err != nil {
				log.Printf("c.consumer.StoreMessage: %v\n", err.Error())
				continue
			}
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}

	return c.consumer.Close()
}
