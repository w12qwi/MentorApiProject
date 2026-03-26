package calculationService

import (
	"MentorApiProject/internal/domain/models"
	"MentorApiProject/internal/infrastructure/grpc"
	"context"
	"errors"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"time"
)

type Repository interface {
	SaveCalculation(ctx context.Context, calculation *models.Calculation) error
	GetCalculation(ctx context.Context, id string) (*models.Calculation, error)
	GetAllCalculations(ctx context.Context) ([]*models.Calculation, error)
	GetCalculationsWithFilters(ctx context.Context, filters models.CalculationsFilters) ([]*models.Calculation, error)
}
type Service struct {
	storage Repository
	tracer  trace.Tracer
}

func NewService(storage Repository, tracer trace.Tracer) *Service {
	return &Service{storage: storage, tracer: tracer}
}

func (s *Service) Calculate(ctx context.Context, calculation models.Calculation) (float64, error) {
	ctx, span := s.tracer.Start(ctx, "service.calculate")
	defer span.End()

	calculation.Id = uuid.New()
	calculation.CreatedAt = time.Now().UTC()

	switch calculation.Sign {
	case "+":
		calculation.Result = calculation.NumA + calculation.NumB
	case "-":
		calculation.Result = calculation.NumA - calculation.NumB
	case "*":
		calculation.Result = calculation.NumA * calculation.NumB
	case "/":
		calculation.Result = calculation.NumA / calculation.NumB
	}

	span.SetAttributes(
		attribute.String("calculation.id", calculation.Id.String()),
		attribute.String("calculation.sign", calculation.Sign),
		attribute.Float64("calculation.num_a", calculation.NumA),
		attribute.Float64("calculation.num_b", calculation.NumB),
		attribute.Float64("calculation.result", calculation.Result),
	)

	err := s.storage.SaveCalculation(ctx, &calculation)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return 0, UnableToSaveCalculationError
	}
	span.SetStatus(codes.Ok, "")

	return calculation.Result, nil
}

func (s *Service) GetCalculation(ctx context.Context, id string) (*models.Calculation, error) {
	ctx, span := s.tracer.Start(ctx, "service.getCalculation")
	defer span.End()
	span.SetAttributes(attribute.String("calculation.id", id))

	result, err := s.storage.GetCalculation(ctx, id)

	if err != nil {
		if errors.Is(err, grpc.CalcultionDoesNotExist) {
			span.SetAttributes(attribute.Bool("calculation.not_found", true))
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return nil, err
		}
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, InternalError
	}
	return result, nil
}

func (s *Service) GetAllCalculations(ctx context.Context) ([]*models.Calculation, error) {
	ctx, span := s.tracer.Start(ctx, "service.getAllCalculations")
	defer span.End()

	result, err := s.storage.GetAllCalculations(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, InternalError
	}

	return result, nil
}

func (s *Service) GetCalculationsWithFilters(ctx context.Context, filters models.CalculationsFilters) ([]*models.Calculation, error) {
	ctx, span := s.tracer.Start(ctx, "service.getCalculationsWithFilters")
	defer span.End()

	result, err := s.storage.GetCalculationsWithFilters(ctx, filters)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return nil, InternalError
	}

	return result, nil
}
