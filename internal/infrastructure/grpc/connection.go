package grpc

import (
	"MentorApiProject/internal/config"
	"context"
	"fmt"
	pb "github.com/w12qwi/calculationsProto/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"log"
	"time"
)

func NewClient(ctx context.Context, cfg config.GRPCconfig) pb.CalculationsDataServiceClient {

	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             3 * time.Second,
		PermitWithoutStream: true,
	}

	dialOpts := []grpc.DialOption{
		grpc.WithKeepaliveParams(kacp),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.DialContext(ctx,
		fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		dialOpts...)

	if err != nil {
		log.Fatal(err)
	}

	client := pb.NewCalculationsDataServiceClient(conn)
	log.Println("Connected to gRPC server")

	return client
}
