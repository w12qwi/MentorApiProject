package models

import (
	"github.com/google/uuid"
	"time"
)

type Calculation struct {
	Id        uuid.UUID
	NumA      int
	NumB      int
	Sign      string
	Result    int
	CreatedAt time.Time
}
