package grpc

import (
	"MentorApiProject/internal/adapter"
	"MentorApiProject/internal/domain/models"
	"context"
	pb "github.com/w12qwi/calculationsProto/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type CalculationsClient struct {
	client  pb.CalculationsDataServiceClient
	timeout time.Duration
}

func NewCalculationsClient(client pb.CalculationsDataServiceClient, timeout int) *CalculationsClient {
	return &CalculationsClient{client: client, timeout: time.Duration(timeout) * time.Second}
}

func (c *CalculationsClient) GetCalculation(ctx context.Context, req *pb.GetCalculationRequest) (*models.Calculation, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.GetCalculation(ctx, req)
	if err != nil {
		status, ok := status.FromError(err)
		if ok {
			switch status.Code() {
			case codes.NotFound:
				return nil, CalcultionDoesNotExist
			}
		}
	}

	calculation, err := adapter.PbCalculationToDomain(resp)
	if err != nil {
		return nil, err
	}

	return calculation, nil
}

func (c *CalculationsClient) GetAllCalculations(ctx context.Context) ([]*models.Calculation, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.GetAllCalculations(ctx, &emptypb.Empty{})
	if err != nil {
		return nil, err
	}

	result, err := adapter.PbCalculationsToDomain(resp.Calculation)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *CalculationsClient) GetCalculationsWithFilter(ctx context.Context, req *pb.GetCalculationsRequest) ([]*models.Calculation, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.GetCalculations(ctx, req)
	if err != nil {
		return nil, err
	}

	result, err := adapter.PbCalculationsToDomain(resp.Calculation)
	if err != nil {
		return nil, err
	}

	return result, nil
}
