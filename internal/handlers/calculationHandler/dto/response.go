package dto

import (
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}
type CalculateResponse struct {
	Result float64 `json:"result"`
}

type CalculationResponse struct {
	Id        string    `json:"id"`
	NumA      float64   `json:"numA"`
	NumB      float64   `json:"numB"`
	Sign      string    `json:"sign"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"createdAt"`
}

type CalculationsResponse struct {
	Calculations []CalculationResponse `json:"calculations"`
}
