package handler

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleMessage(message []byte, offset kafka.Offset) error {
	log.Printf("Message from kafka with offset %d '%s", offset, message)
	return nil
}
