package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

func NewConsumer(broker, topic, group string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: group,
	})
}

func Consume(reader *kafka.Reader, handler func([]byte)) {
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			continue
		}
		handler(msg.Value)
	}
}
