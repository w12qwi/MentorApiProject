package calculationService

import (
	"MentorApiProject/internal/domain/models"
	"MentorApiProject/internal/infrastructure/grpc"
	"context"
	"errors"
	"github.com/google/uuid"
	"log/slog"
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
}

func NewService(storage Repository) *Service {
	return &Service{storage: storage}
}

func (s *Service) Calculate(ctx context.Context, calculation models.Calculation) (float64, error) {

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

	err := s.storage.SaveCalculation(ctx, &calculation)
	if err != nil {
		return 0, UnableToSaveCalculationError
	}

	return calculation.Result, nil
}

func (s *Service) GetCalculation(ctx context.Context, id string) (*models.Calculation, error) {
	result, err := s.storage.GetCalculation(ctx, id)

	if err != nil {
		if errors.Is(err, grpc.CalcultionDoesNotExist) {
			return nil, err
		}
		slog.Error("Unable to get calculation by id: ", err)
		return nil, InternalError
	}
	return result, nil
}

func (s *Service) GetAllCalculations(ctx context.Context) ([]*models.Calculation, error) {
	result, err := s.storage.GetAllCalculations(ctx)
	if err != nil {
		slog.Error("Unable to get all calculations: ", err)
		return nil, InternalError
	}

	return result, nil
}

func (s *Service) GetCalculationsWithFilters(ctx context.Context, filters models.CalculationsFilters) ([]*models.Calculation, error) {
	result, err := s.storage.GetCalculationsWithFilters(ctx, filters)
	if err != nil {
		slog.Error("Unable to get calculations with filters: ", err)
		return nil, InternalError
	}

	return result, nil
}
