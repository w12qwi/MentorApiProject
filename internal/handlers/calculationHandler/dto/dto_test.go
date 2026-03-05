package dto

import (
	"github.com/google/uuid"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateRequest_Validate(t *testing.T) {

	tests := []struct {
		name    string
		req     CalculateRequest
		wantErr error
	}{
		{
			name: "Valid addition",
			req: CalculateRequest{
				NumA: 10,
				NumB: 20,
				Sign: "+",
			},
			wantErr: nil,
		}, {
			name: "Valid subtraction",
			req: CalculateRequest{
				NumA: 10,
				NumB: 20,
				Sign: "-",
			},
			wantErr: nil,
		}, {
			name: "Valid multiplication",
			req: CalculateRequest{
				NumA: 10,
				NumB: 20,
				Sign: "*",
			},
			wantErr: nil,
		}, {
			name: "Valid division",
			req: CalculateRequest{
				NumA: 10,
				NumB: 20,
				Sign: "/",
			},
			wantErr: nil,
		}, {
			name: "Invalid sign",
			req: CalculateRequest{
				NumA: 10,
				NumB: 20,
				Sign: "?",
			},
			wantErr: InvalidSignError,
		}, {
			name: "Division by zero",
			req: CalculateRequest{
				NumA: 10,
				NumB: 0,
				Sign: "/",
			},
			wantErr: DivisionByZeroError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetByIdRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     GetByIdRequest
		wantErr error
	}{
		{
			name: "Valid UUID",
			req: GetByIdRequest{
				Id: uuid.New().String(),
			},
			wantErr: nil,
		}, {
			name: "Invalid UUID",
			req: GetByIdRequest{
				Id: "invalid-uuid",
			},
			wantErr: InvalidIdFormatError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetCalculationsByDateRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     GetCalculationsByDateRequest
		wantErr error
	}{
		{
			name: "Valid date",
			req: GetCalculationsByDateRequest{
				Date: "2023-01-01",
			},
			wantErr: nil,
		}, {
			name: "Invalid date ",
			req: GetCalculationsByDateRequest{
				Date: "0101202305005",
			},
			wantErr: InvalidDateFormatError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestGetCalculationsByDateRangeRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     GetCalculationsByDateRangeRequest
		wantErr error
	}{
		{
			name: "Valid date range",
			req: GetCalculationsByDateRangeRequest{
				From: "2023-01-01",
				To:   "2023-01-02",
			},
			wantErr: nil,
		}, {
			name: "Invalid date range",
			req: GetCalculationsByDateRangeRequest{
				From: "2023-01-02",
				To:   "2023-01-01",
			},
			wantErr: InvalidDateRangeError,
		}, {
			name: "Invalid date format",
			req: GetCalculationsByDateRangeRequest{
				From: "0101202305005",
				To:   "0101202305005",
			},
			wantErr: InvalidDateFormatError,
		}, {
			name: "Valid date format",
			req: GetCalculationsByDateRangeRequest{
				From: "2023-01-01",
				To:   "2023-01-02",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
