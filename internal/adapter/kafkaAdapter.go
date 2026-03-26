package adapter

import (
	"MentorApiProject/internal/domain/models"
	"MentorApiProject/internal/infrastructure/kafka"
)

func CalculationToKafkaMessage(calculation *models.Calculation) *kafka.CalculationMessage {
	return &kafka.CalculationMessage{
		Id:        calculation.Id.String(),
		NumA:      calculation.NumA,
		NumB:      calculation.NumB,
		Sign:      calculation.Sign,
		Result:    calculation.Result,
		CreatedAt: calculation.CreatedAt,
	}
}
