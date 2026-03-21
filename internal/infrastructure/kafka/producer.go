package kafka

import (
	"MentorApiProject/internal/config"
	"context"
	"encoding/json"
	"fmt"
	k "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *k.Writer
}

func NewProducer(cfg config.KafkaConfig) *Producer {

	return &Producer{
		writer: &k.Writer{
			Addr:         k.TCP(fmt.Sprintf("%s:%s", cfg.BrokerHost, cfg.BrokerPort)),
			Topic:        cfg.Topic,
			Balancer:     &k.LeastBytes{},
			RequiredAcks: k.RequireOne,
			Compression:  k.Snappy,
		},
	}
}

func (p *Producer) Publish(ctx context.Context, calculation CalculationMessage) (err error) {

	msg, err := json.Marshal(calculation)
	if err != nil {
		return err
	}

	err = p.writer.WriteMessages(ctx, k.Message{
		Key:   []byte(calculation.Id),
		Value: msg})
	if err != nil {
		return err
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
