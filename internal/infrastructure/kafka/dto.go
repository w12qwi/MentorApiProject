package kafka

import "time"

type CalculationMessage struct {
	Id        string    `json:"id"`
	NumA      float64   `json:"numA"`
	NumB      float64   `json:"numB"`
	Sign      string    `json:"sign"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"createdAt"`
}
