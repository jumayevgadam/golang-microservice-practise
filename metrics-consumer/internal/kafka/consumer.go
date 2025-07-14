package kafka

import (
	"fmt"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	sessionTimeoutMs = 7000
	noTimeout        = -1
)

type Consumer struct {
	consumer *kafka.Consumer
	stop     bool
}

func NewConsumer(address []string, topic, consumerGroup string) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":        strings.Join(address, ","),
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

	return &Consumer{consumer: c}, nil
}
