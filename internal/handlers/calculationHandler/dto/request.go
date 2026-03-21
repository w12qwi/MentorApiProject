package dto

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

var (
	DivisionByZeroError        = errors.New("Cannot divide by zero")
	InvalidSignError           = errors.New("Invalid sign")
	InvalidIdFormatError       = errors.New("Invalid id format")
	InvalidDateFormatError     = errors.New("Invalid date format")
	InvalidRequestBodyError    = errors.New("Invalid request body")
	InvalidDateRangeError      = errors.New("Invalid date range")
	InvalidDateConstraintError = errors.New("Invalid date constraint")
)

type CalculateRequest struct {
	NumA float64 `json:"numA"`
	NumB float64 `json:"numB"`
	Sign string  `json:"sign"`
}

type GetCalculationstWithFiltersRequest struct {
	Sign     *string `json:"sign"`
	Date     *string `json:"date"`
	DateFrom *string `json:"dateFrom"`
	DateTo   *string `json:"dateTo"`
}

type GetCalculationByIdRequest struct {
	Id string `json:"id"`
}

func (r *GetCalculationstWithFiltersRequest) Validate() error {
	if r.Date != nil && r.DateFrom != nil && r.DateTo != nil || r.Date == nil && r.DateFrom == nil || r.Date != nil && r.DateTo == nil {
		return InvalidDateConstraintError
	}

	if !strings.Contains("+-*/", *r.Sign) && r.Sign != nil {
		return InvalidSignError
	}

	_, err := time.Parse(dateFormat, *r.Date)
	if err != nil {
		return InvalidDateFormatError
	}

	from, err := time.Parse(dateFormat, *r.DateFrom)
	if err != nil {
		return InvalidDateFormatError
	}

	to, err := time.Parse(dateFormat, *r.DateTo)
	if err != nil {
		return InvalidDateFormatError
	}

	if from.After(to) {
		return InvalidDateRangeError
	}

	return nil
}

func (r *CalculateRequest) Validate() error {
	if r.NumB == 0 && r.Sign == "/" {
		return DivisionByZeroError
	}

	if !strings.Contains("+-*/", r.Sign) {
		return InvalidSignError
	}

	return nil

}

func (r *GetCalculationByIdRequest) Validate() error {
	err := uuid.Validate(r.Id)
	if err != nil {
		return InvalidIdFormatError
	}

	return nil
}
