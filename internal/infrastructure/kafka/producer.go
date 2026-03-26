package kafka

import (
	"MentorApiProject/internal/config"
	"context"
	"encoding/json"
	"fmt"
	k "github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"log"
	"time"
)

type Producer struct {
	writer *k.Writer
	tracer trace.Tracer
}

func NewProducer(cfg config.KafkaConfig, tracer trace.Tracer) *Producer {

	return &Producer{
		writer: &k.Writer{
			Addr:         k.TCP(fmt.Sprintf("%s:%s", cfg.BrokerHost, cfg.BrokerPort)),
			Topic:        cfg.Topic,
			Balancer:     &k.LeastBytes{},
			RequiredAcks: k.RequireOne,
			Compression:  k.Snappy,
			BatchTimeout: 10 * time.Millisecond,
		},
		tracer: tracer,
	}
}

func (p *Producer) Produce(ctx context.Context, calculation CalculationMessage) error {
	ctx, span := p.tracer.Start(ctx, "kafka.produce")
	defer span.End()

	span.SetAttributes(
		attribute.String("messaging.system", "kafka"),
		attribute.String("messaging.destination", p.writer.Topic),
		attribute.String("messaging.destination_kind", "topic"),
		attribute.String("messaging.operation", "publish"),
		attribute.String("messaging.kafka.message_key", calculation.Id),
	)

	msg, err := json.Marshal(calculation)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "marshal failed")
		return err
	}

	carrier := propagation.MapCarrier{}
	otel.GetTextMapPropagator().Inject(ctx, carrier)

	headers := make([]k.Header, 0, len(carrier))
	for key, val := range carrier {
		headers = append(headers, k.Header{
			Key:   key,
			Value: []byte(val),
		})
	}

	log.Println("OUTGOING HEADERS:")
	for key, val := range carrier {
		log.Printf("%s = %s\n", key, val)
	}

	start := time.Now()
	err = p.writer.WriteMessages(ctx, k.Message{
		Key:     []byte(calculation.Id),
		Value:   msg,
		Headers: headers,
	})
	log.Printf("WriteMessages took: %s", time.Since(start))
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "write kafka message failed")
		return err
	}

	span.SetStatus(codes.Ok, "")
	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
