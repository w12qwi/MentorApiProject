package repository

import (
	"MentorApiProject/internal/adapter"
	"MentorApiProject/internal/domain/models"
	"MentorApiProject/internal/infrastructure/kafka"
	"context"
	pb "github.com/w12qwi/calculationsProto/gen"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type GRPCClient interface {
	GetCalculation(ctx context.Context, req *pb.GetCalculationRequest) (*models.Calculation, error)
	GetAllCalculations(ctx context.Context) ([]*models.Calculation, error)
	GetCalculations(ctx context.Context, req *pb.GetCalculationsRequest) ([]*models.Calculation, error)
}

type KafkaProducer interface {
	Produce(ctx context.Context, calculation kafka.CalculationMessage) (err error)
}

type CalculationsRepository struct {
	grpcClient    GRPCClient
	kafkaProducer KafkaProducer
	tracer        trace.Tracer
}

func NewCalculationsRepository(grpcClient GRPCClient, kafkaProducer KafkaProducer, tracer trace.Tracer) *CalculationsRepository {
	return &CalculationsRepository{grpcClient: grpcClient, kafkaProducer: kafkaProducer, tracer: tracer}
}

func (r *CalculationsRepository) SaveCalculation(ctx context.Context, calculation *models.Calculation) error {
	ctx, span := r.tracer.Start(ctx, "repository.saveCalculation")
	defer span.End()

	span.SetAttributes(attribute.String("repository.action", "enqueue_calculation"))

	msg := adapter.CalculationToKafkaMessage(calculation)

	err := r.kafkaProducer.Produce(ctx, *msg)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "write kafka message failed")
		return err
	}

	span.SetStatus(codes.Ok, "")

	return nil
}

func (r *CalculationsRepository) GetCalculation(ctx context.Context, id string) (*models.Calculation, error) {
	ctx, span := r.tracer.Start(ctx, "repository.getCalculation")
	defer span.End()
	span.SetAttributes(attribute.String("repository.action", "get_calculation"),
		attribute.String("calculation.id", id),
		attribute.String("repository.grpc_method", "GetCalculation"),
		attribute.String("repository.grpc_service", "CalculationsDataService"),
	)

	request := pb.GetCalculationRequest{Id: id}
	resp, err := r.grpcClient.GetCalculation(ctx, &request)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get calculation failed")
		return nil, err
	}
	span.SetStatus(codes.Ok, "")

	return resp, nil
}

func (r *CalculationsRepository) GetAllCalculations(ctx context.Context) ([]*models.Calculation, error) {
	ctx, span := r.tracer.Start(ctx, "repository.getAllCalculations")
	defer span.End()

	resp, err := r.grpcClient.GetAllCalculations(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get all calculations failed")
		return nil, err
	}
	return resp, nil
}

func (r *CalculationsRepository) GetCalculationsWithFilters(ctx context.Context, filters models.CalculationsFilters) ([]*models.Calculation, error) {
	ctx, span := r.tracer.Start(ctx, "repository.getCalculationsWithFilters")
	defer span.End()

	req, err := adapter.DomainFiltersToPb(&filters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to convert domain filters to proto filters")
		return nil, err
	}
	resp, err := r.grpcClient.GetCalculations(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "get calculations with filters failed")
		return nil, err
	}
	return resp, nil
}
