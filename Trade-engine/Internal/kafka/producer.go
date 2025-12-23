package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:  kafka.TCP(broker),
			Topic: topic,
		},
	}
}

func (p *Producer) Publish(msg []byte) error {
	return p.writer.WriteMessages(context.Background(),
		kafka.Message{Value: msg},
	)
}
