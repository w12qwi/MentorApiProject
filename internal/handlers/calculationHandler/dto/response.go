package dto

import (
	"time"
)

type ErrorResponse struct {
	Error string `json:"error"`
}
type CalculateResponse struct {
	Result int `json:"result"`
}

type CalculationResponse struct {
	Id        string    `json:"id"`
	NumA      int       `json:"numA"`
	NumB      int       `json:"numB"`
	Sign      string    `json:"sign"`
	Result    int       `json:"result"`
	CreatedAt time.Time `json:"createdAt"`
}

type CalculationsResponse struct {
	Calculations []CalculationResponse `json:"calculations"`
}
