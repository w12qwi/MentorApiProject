package repository

import (
	"MentorApiProject/internal/adapter"
	"MentorApiProject/internal/domain/models"
	"MentorApiProject/internal/infrastructure/kafka"
	"context"
	pb "github.com/w12qwi/calculationsProto/gen"
)

type GRPCClient interface {
	GetCalculation(ctx context.Context, req *pb.GetCalculationRequest) (*models.Calculation, error)
	GetAllCalculations(ctx context.Context) ([]*models.Calculation, error)
	GetCalculationsWithFilter(ctx context.Context, req *pb.GetCalculationsRequest) ([]*models.Calculation, error)
}

type KafkaProducer interface {
	Publish(ctx context.Context, calculation kafka.CalculationMessage) (err error)
}

type CalculationsRepository struct {
	grpcClient    GRPCClient
	kafkaProducer KafkaProducer
}

func NewCalculationsRepository(grpcClient GRPCClient, kafkaProducer KafkaProducer) *CalculationsRepository {
	return &CalculationsRepository{grpcClient: grpcClient, kafkaProducer: kafkaProducer}
}

func (r *CalculationsRepository) SaveCalculation(ctx context.Context, calculation *models.Calculation) error {
	msg, err := adapter.CalculationToKafkaMessage(calculation)
	if err != nil {
		return err
	}
	err = r.kafkaProducer.Publish(ctx, *msg)
	if err != nil {
		return err
	}
	return nil
}

func (r *CalculationsRepository) GetCalculation(ctx context.Context, id string) (*models.Calculation, error) {
	request := pb.GetCalculationRequest{Id: id}
	resp, err := r.grpcClient.GetCalculation(ctx, &request)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *CalculationsRepository) GetAllCalculations(ctx context.Context) ([]*models.Calculation, error) {
	resp, err := r.grpcClient.GetAllCalculations(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *CalculationsRepository) GetCalculationsWithFilters(ctx context.Context, filters models.CalculationsFilters) ([]*models.Calculation, error) {
	req, err := adapter.DomainFiltersToPb(&filters)
	if err != nil {
		return nil, err
	}
	resp, err := r.grpcClient.GetCalculationsWithFilter(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
