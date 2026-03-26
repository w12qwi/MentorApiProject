package kafka

import "time"

type CalculationMessage struct {
	Id        string    `json:"id"`
	NumA      float64   `json:"num_a"`
	NumB      float64   `json:"num_b"`
	Sign      string    `json:"sign"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}
