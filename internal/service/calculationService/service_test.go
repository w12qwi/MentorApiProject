package calculationService

import (
	"MentorApiProject/internal/domain/models"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type mockStorage struct {
	saveErr        error
	getByIdResp    models.Calculation
	getByIdErr     error
	getAllResp     []models.Calculation
	getAllErr      error
	getByDateResp  []models.Calculation
	getByDateErr   error
	getByRangeResp []models.Calculation
	getByRangeErr  error
}

func (m *mockStorage) SaveCalculation(_ context.Context, _ models.Calculation) error {
	return m.saveErr
}

func (m *mockStorage) GetById(_ context.Context, _ uuid.UUID) (models.Calculation, error) {
	return m.getByIdResp, m.getByIdErr
}

func (m *mockStorage) GetAllCalculations(_ context.Context) ([]models.Calculation, error) {
	return m.getAllResp, m.getAllErr
}

func (m *mockStorage) GetCalculationsByDate(_ context.Context, _ time.Time) ([]models.Calculation, error) {
	return m.getByDateResp, m.getByDateErr
}

func (m *mockStorage) GetCalculationsByDateRange(_ context.Context, _, _ time.Time) ([]models.Calculation, error) {
	return m.getByRangeResp, m.getByRangeErr
}

func TestService_Calculate(t *testing.T) {
	tests := []struct {
		name       string
		input      models.Calculation
		saveErr    error
		wantResult float64
		wantErr    error
	}{
		{
			name:       "addition",
			input:      models.Calculation{NumA: 2, NumB: 3, Sign: "+"},
			wantResult: 5,
		},
		{
			name:       "subtraction",
			input:      models.Calculation{NumA: 5, NumB: 3, Sign: "-"},
			wantResult: 2,
		},
		{
			name:       "multiplication",
			input:      models.Calculation{NumA: 2, NumB: 3, Sign: "*"},
			wantResult: 6,
		},
		{
			name:       "division",
			input:      models.Calculation{NumA: 6, NumB: 3, Sign: "/"},
			wantResult: 2,
		},
		{
			name:    "storage error",
			input:   models.Calculation{NumA: 2, NumB: 3, Sign: "+"},
			saveErr: sql.ErrConnDone,
			wantErr: UnableToSaveCalculationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(&mockStorage{saveErr: tt.saveErr})
			result, err := service.Calculate(context.Background(), tt.input)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantResult, result)
			}
		})
	}
}

func TestService_GetById(t *testing.T) {
	id := uuid.New()
	calculation := models.Calculation{Id: id, NumA: 2, NumB: 3, Sign: "+", Result: 5}

	tests := []struct {
		name     string
		mock     mockStorage
		wantResp models.Calculation
		wantErr  error
	}{
		{
			name:     "found",
			mock:     mockStorage{getByIdResp: calculation},
			wantResp: calculation,
		},
		{
			name:    "not found",
			mock:    mockStorage{getByIdErr: sql.ErrNoRows},
			wantErr: NoSuchCalculationError,
		},
		{
			name:    "internal error",
			mock:    mockStorage{getByIdErr: sql.ErrConnDone},
			wantErr: InternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(&tt.mock)
			resp, err := service.GetById(context.Background(), id)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestService_GetAllCalculations(t *testing.T) {
	calculations := []models.Calculation{
		{Id: uuid.New(), NumA: 1, NumB: 2, Sign: "+", Result: 3},
		{Id: uuid.New(), NumA: 5, NumB: 3, Sign: "-", Result: 2},
	}

	tests := []struct {
		name     string
		mock     mockStorage
		wantResp []models.Calculation
		wantErr  error
	}{
		{
			name:     "success",
			mock:     mockStorage{getAllResp: calculations},
			wantResp: calculations,
		},
		{
			name:    "internal error",
			mock:    mockStorage{getAllErr: sql.ErrConnDone},
			wantErr: InternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(&tt.mock)
			resp, err := service.GetAllCalculations(context.Background())
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestService_GetCalculationsByDate(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	calculations := []models.Calculation{
		{Id: uuid.New(), NumA: 1, NumB: 2, Sign: "+", Result: 3, CreatedAt: date},
	}

	tests := []struct {
		name     string
		mock     mockStorage
		wantResp []models.Calculation
		wantErr  error
	}{
		{
			name:     "success",
			mock:     mockStorage{getByDateResp: calculations},
			wantResp: calculations,
		},
		{
			name:    "internal error",
			mock:    mockStorage{getByDateErr: sql.ErrConnDone},
			wantErr: InternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(&tt.mock)
			resp, err := service.GetCalculationsByDate(context.Background(), date)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}

func TestService_GetCalculationsByDateRange(t *testing.T) {
	from := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC)
	calculations := []models.Calculation{
		{Id: uuid.New(), NumA: 1, NumB: 2, Sign: "+", Result: 3, CreatedAt: from},
		{Id: uuid.New(), NumA: 5, NumB: 3, Sign: "-", Result: 2, CreatedAt: to},
	}

	tests := []struct {
		name     string
		mock     mockStorage
		wantResp []models.Calculation
		wantErr  error
	}{
		{
			name:     "success",
			mock:     mockStorage{getByRangeResp: calculations},
			wantResp: calculations,
		},
		{
			name:    "internal error",
			mock:    mockStorage{getByRangeErr: sql.ErrConnDone},
			wantErr: InternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(&tt.mock)
			resp, err := service.GetCalculationsByDateRange(context.Background(), from, to)
			assert.ErrorIs(t, err, tt.wantErr)
			if tt.wantErr == nil {
				assert.Equal(t, tt.wantResp, resp)
			}
		})
	}
}
