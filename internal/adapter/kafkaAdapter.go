package adapter

import (
	"MentorApiProject/internal/domain/models"
	"MentorApiProject/internal/infrastructure/kafka"
)

func CalculationToKafkaMessage(calculation *models.Calculation) (*kafka.CalculationMessage, error) {

}
