package grpc

import (
	"MentorApiProject/internal/adapter"
	"MentorApiProject/internal/domain/models"
	"context"
	"fmt"
	pb "github.com/w12qwi/calculationsProto/gen"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

type CalculationsClient struct {
	client  pb.CalculationsDataServiceClient
	timeout time.Duration
	tracer  trace.Tracer
}

func NewCalculationsClient(client pb.CalculationsDataServiceClient, timeout int, tracer trace.Tracer) *CalculationsClient {
	return &CalculationsClient{client: client, timeout: time.Duration(timeout) * time.Second, tracer: tracer}
}

func (c *CalculationsClient) GetCalculation(ctx context.Context, req *pb.GetCalculationRequest) (*models.Calculation, error) {
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	resp, err := c.client.GetCalculation(ctx, req)
	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, CalcultionDoesNotExist
		}
		return nil, err
	}

	if resp == nil {
		return nil, fmt.Errorf("empty response from gRPC server for id=%s", req.Id)
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

func (c *CalculationsClient) GetCalculations(ctx context.Context, req *pb.GetCalculationsRequest) ([]*models.Calculation, error) {
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
