package util

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func PublishKafkaMessage(writer *kafka.Writer, key string, message []byte) error {
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(key),
		Value: []byte(message),
	})

	return err
}
