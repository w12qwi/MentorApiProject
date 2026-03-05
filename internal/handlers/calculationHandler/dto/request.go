package dto

import (
	"errors"
	"github.com/google/uuid"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

var (
	DivisionByZeroError     = errors.New("Cannot divide by zero")
	InvalidSignError        = errors.New("Invalid sign")
	InvalidIdFormatError    = errors.New("Invalid id format")
	InvalidDateFormatError  = errors.New("Invalid date format")
	InvalidRequestBodyError = errors.New("Invalid request body")
	InvalidDateRangeError   = errors.New("Invalid date range")
)

type CalculateRequest struct {
	NumA float64 `json:"numA"`
	NumB float64 `json:"numB"`
	Sign string  `json:"sign"`
}

type GetByIdRequest struct {
	Id string `json:"id"`
}
type GetCalculationsByDateRequest struct {
	Date string `json:"date"`
}

type GetCalculationsByDateRangeRequest struct {
	From string `json:"fromDate"`
	To   string `json:"toDate"`
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

func (r *GetByIdRequest) Validate() error {
	err := uuid.Validate(r.Id)
	if err != nil {
		return InvalidIdFormatError
	}

	return nil
}

func (r *GetCalculationsByDateRequest) Validate() error {
	_, err := time.Parse(dateFormat, r.Date)
	if err != nil {
		return InvalidDateFormatError
	}

	return nil
}

func (r *GetCalculationsByDateRangeRequest) Validate() error {
	from, err := time.Parse(dateFormat, r.From)
	if err != nil {
		return InvalidDateFormatError
	}

	to, err := time.Parse(dateFormat, r.To)
	if err != nil {
		return InvalidDateFormatError
	}

	if from.After(to) {
		return InvalidDateRangeError
	}

	return nil
}

func (r *GetByIdRequest) UUID() uuid.UUID {
	return uuid.MustParse(r.Id)
}

func (r *GetCalculationsByDateRequest) UTCDate() time.Time {
	date, _ := time.Parse(dateFormat, r.Date)
	return date.UTC()
}

func (r *GetCalculationsByDateRangeRequest) UTCDateFrom() time.Time {
	from, _ := time.Parse(dateFormat, r.From)

	return from.UTC()
}

func (r *GetCalculationsByDateRangeRequest) UTCDateTo() time.Time {
	to, _ := time.Parse(dateFormat, r.To)

	return to.UTC().Add(24*time.Hour - time.Second)
}
