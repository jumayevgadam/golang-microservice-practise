package internal

import (
	"context"
	"log"
	"metrics-consumer/internal/handler"
	"metrics-consumer/internal/kafka"
	"os"
	"os/signal"
	"syscall"
)

func BootStrapMetricsService(envPath string) error {
	if err := LoadEnv(envPath); err != nil {
		return err
	}

	cfg, err := NewMetricsConsumerConfig()
	if err != nil {
		return err
	}

	h := handler.NewHandler()

	consumer, err := kafka.NewConsumer(
		h,
		cfg.Brokers,
		cfg.Topic,
		cfg.ConsumerGroup,
	)
	if err != nil {
		log.Printf("failed to create metrics consumer service: %v\n", err.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		consumer.Start(ctx)
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan

	if err := consumer.Stop(); err != nil {
		log.Printf("error shutting down metrics-consumer service: %v\n", err.Error())
	}

	return nil
}
