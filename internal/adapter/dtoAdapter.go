package adapter

import (
	"MentorApiProject/internal/domain/models"
	"MentorApiProject/internal/handlers/calculationHandler/dto"
	"time"
)

const dateFormat = "2006-01-02"

func DtoFiltersToDomain(req dto.GetCalculationstWithFiltersRequest) (models.CalculationsFilters, error) {
	var result models.CalculationsFilters

	if req.Sign != nil {
		result.Sign = req.Sign
	}

	if req.Date != nil {
		parsedDate, err := time.Parse(dateFormat, *req.Date)
		if err != nil {
			return models.CalculationsFilters{}, err
		}
		result.Date = &parsedDate
	}

	if req.DateFrom != nil {
		parsedDateFrom, err := time.Parse(dateFormat, *req.DateFrom)
		if err != nil {
			return models.CalculationsFilters{}, err
		}
		result.DateFrom = &parsedDateFrom
	}

	if req.DateTo != nil {
		parsedDateTo, err := time.Parse(dateFormat, *req.DateTo)
		if err != nil {
			return models.CalculationsFilters{}, err
		}
		result.DateTo = &parsedDateTo
	}

	return result, nil
}

func DomainToDto(calculation *models.Calculation) dto.CalculationResponse {
	return dto.CalculationResponse{
		Id:        calculation.Id.String(),
		NumA:      calculation.NumA,
		NumB:      calculation.NumB,
		Sign:      calculation.Sign,
		Result:    calculation.Result,
		CreatedAt: calculation.CreatedAt,
	}
}

func DomainSliceToDto(calculations []*models.Calculation) []dto.CalculationResponse {
	result := make([]dto.CalculationResponse, 0)
	for _, calculation := range calculations {
		result = append(result, DomainToDto(calculation))
	}
	return result
}
