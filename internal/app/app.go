package app

import (
	"MentorApiProject/internal/config"
	"MentorApiProject/internal/handlers/calculationHandler"

	"MentorApiProject/internal/infrastructure/grpc"
	"MentorApiProject/internal/infrastructure/kafka"
	"MentorApiProject/internal/repository"
	"MentorApiProject/internal/service/calculationService"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	server *http.Server
}

func NewApp() *App {

	grpcConfig := config.LoadGRPCConfig()
	kafkaConfig := config.LoadKafkaConfig()
	AppConfig := config.LoadAppConfig()

	grpcClient := grpc.NewClient(context.TODO(), grpcConfig)
	calculationGRPCClient := grpc.NewCalculationsClient(grpcClient, 10)

	kafkaProducer := kafka.NewProducer(kafkaConfig)

	repo := repository.NewCalculationsRepository(calculationGRPCClient, kafkaProducer)

	calculationService := calculationService.NewService(repo)

	HTTPHandler := calculationHandler.NewHandler(calculationService)

	mux := http.NewServeMux()
	HTTPHandler.RegisterRoutes(mux)

	return &App{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%s", AppConfig.Port),
			Handler: mux,
		},
	}
}

func (a *App) Run() error {
	slog.Info("Starting server on port 8080")
	return a.server.ListenAndServe()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	return err
}
