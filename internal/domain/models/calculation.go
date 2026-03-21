package models

import (
	"github.com/google/uuid"
	"time"
)

type Calculation struct {
	Id        uuid.UUID
	NumA      float64
	NumB      float64
	Sign      string
	Result    float64
	CreatedAt time.Time
}

type CalculationsFilters struct {
	Sign     *string
	Date     *time.Time
	DateFrom *time.Time
	DateTo   *time.Time
}
