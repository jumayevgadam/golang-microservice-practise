package internal

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config represent metrics-consumer kafka service configurations.
type Config struct {
	Brokers       string
	Topic         string
	ConsumerGroup string
}

func LoadEnv(path string) error {
	if err := godotenv.Load(path); err != nil {
		return fmt.Errorf("failed to load environment variables: %w", err)
	}

	return nil
}

func NewMetricsConsumerConfig() (*Config, error) {
	cfg := &Config{
		Brokers:       os.Getenv("BROKERS"),
		Topic:         os.Getenv("TOPIC"),
		ConsumerGroup: os.Getenv("CONSUMER_GROUP"),
	}

	return cfg, nil
}
