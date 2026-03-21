package adapter

import (
	"MentorApiProject/internal/domain/models"
	"github.com/google/uuid"
	pb "github.com/w12qwi/calculationsProto/gen"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PbCalculationToDomain(c *pb.Calculation) (*models.Calculation, error) {
	id, err := uuid.Parse(c.Id)
	if err != nil {
		return nil, err
	}

	return &models.Calculation{
		Id:        id,
		NumA:      c.NumA,
		NumB:      c.NumB,
		Sign:      c.Sign,
		Result:    c.Result,
		CreatedAt: c.CreatedAt.AsTime(),
	}, nil
}

func PbCalculationsToDomain(c []*pb.Calculation) ([]*models.Calculation, error) {
	result := make([]*models.Calculation, 0)
	for _, calculation := range c {
		calculation, err := PbCalculationToDomain(calculation)
		if err != nil {
			return nil, err
		}
		result = append(result, calculation)
	}
	return result, nil
}

func DomainFiltersToPb(filters *models.CalculationsFilters) (*pb.GetCalculationsRequest, error) {
	req := &pb.GetCalculationsRequest{
		Sign: filters.Sign,
	}
	if filters.Date != nil {
		req.Date = timestamppb.New(*filters.Date)
	}
	if filters.DateFrom != nil {
		req.DateFrom = timestamppb.New(*filters.DateFrom)
	}
	if filters.DateTo != nil {
		req.DateTo = timestamppb.New(*filters.DateTo)
	}
	return req, nil
}
