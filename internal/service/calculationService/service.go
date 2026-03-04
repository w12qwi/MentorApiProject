package calculationService

import (
	"MentorApiProject/internal/domain/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Storage interface {
	SaveCalculation(ctx context.Context, calculation models.Calculation) error
	GetById(ctx context.Context, id uuid.UUID) (models.Calculation, error)
	GetAllCalculations(ctx context.Context) ([]models.Calculation, error)
	GetCalculationsByDate(ctx context.Context, date time.Time) ([]models.Calculation, error)
	GetCalculationsByDateRange(ctx context.Context, from, to time.Time) ([]models.Calculation, error)
}
type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) Calculate(ctx context.Context, calculation models.Calculation) (int, error) {

	calculation.Id = uuid.New()
	calculation.CreatedAt = time.Now().UTC()

	switch calculation.Sign {
	case "+":
		calculation.Result = calculation.NumA + calculation.NumB
	case "-":
		calculation.Result = calculation.NumA - calculation.NumB
	case "*":
		calculation.Result = calculation.NumA * calculation.NumB //слишком большие значения могут превысить максимальное значение инта
	case "/":
		calculation.Result = calculation.NumA / calculation.NumB
	}

	err := s.storage.SaveCalculation(ctx, calculation)
	if err != nil {
		slog.Error(fmt.Sprintf(UnableToSaveCalculationError.Error())+":%s", err.Error())
		return 0, UnableToSaveCalculationError
	}

	return calculation.Result, nil
}

func (s *Service) GetById(ctx context.Context, id uuid.UUID) (models.Calculation, error) {

	response, err := s.storage.GetById(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Calculation{}, NoSuchCalculationError
	} else if err != nil {
		slog.Error(fmt.Sprintf("Some error occured while querying postgres"+":%s", err.Error()))
		return models.Calculation{}, InternalError
	}

	return response, err
}

func (s *Service) GetAllCalculations(ctx context.Context) ([]models.Calculation, error) {

	result, err := s.storage.GetAllCalculations(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("Some error occured while querying postgres"+":%s", err.Error()))
		return nil, InternalError
	}

	return result, nil
}

func (s *Service) GetCalculationsByDate(ctx context.Context, date time.Time) ([]models.Calculation, error) {

	response, err := s.storage.GetCalculationsByDate(ctx, date)
	if err != nil {
		slog.Error(fmt.Sprintf("Some error occured while querying postgres"+":%s", err.Error()))
		return nil, InternalError
	}

	return response, nil
}

func (s *Service) GetCalculationsByDateRange(ctx context.Context, from, to time.Time) ([]models.Calculation, error) {

	response, err := s.storage.GetCalculationsByDateRange(ctx, from, to)
	if err != nil {
		slog.Error(fmt.Sprintf("Some error occured while querying postgres"+":%s", err.Error()))
		return nil, InternalError
	}

	return response, nil
}
