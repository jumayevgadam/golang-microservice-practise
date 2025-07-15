package internal

import (
	"log"
	"os"
)

func BootStrapMetricsService() error {
	brokers := os.Getenv("KAFKA_BROKERS")
	if brokers == "" {
		log.Println("kafka brokers not set")
	}

	return nil
}
