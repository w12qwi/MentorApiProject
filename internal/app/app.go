package app

import (
	"MentorApiProject/internal/config"
	"MentorApiProject/internal/handlers/calculationHandler"
	"MentorApiProject/internal/infrastructure/grpc"
	"MentorApiProject/internal/infrastructure/kafka"
	"MentorApiProject/internal/infrastructure/tracing"
	"MentorApiProject/internal/repository"
	"MentorApiProject/internal/service/calculationService"
	"context"
	"go.opentelemetry.io/otel"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	server         *http.Server
	tracerShutdown func(ctx context.Context) error
}

func NewApp() *App {

	grpcConfig := config.LoadGRPCConfig()
	kafkaConfig := config.LoadKafkaConfig()
	jaegerConfig := config.LoadJaegerConfig()

	tracerShutdown := tracing.NewTracer("api-service", jaegerConfig)
	tracer := otel.Tracer("api-service")

	grpcClient := grpc.NewClient(context.Background(), grpcConfig)
	calculationGRPCClient := grpc.NewCalculationsClient(grpcClient, grpcConfig.Timeout, tracer)

	kafkaProducer := kafka.NewProducer(kafkaConfig, tracer)

	repo := repository.NewCalculationsRepository(calculationGRPCClient, kafkaProducer, tracer)

	calculationService := calculationService.NewService(repo, tracer)

	HTTPHandler := calculationHandler.NewHandler(calculationService, tracer)

	mux := http.NewServeMux()
	HTTPHandler.RegisterRoutes(mux)

	return &App{
		server: &http.Server{
			Addr:    ":9999",
			Handler: mux,
		},
		tracerShutdown: tracerShutdown,
	}
}

func (a *App) Run() error {
	slog.Info("Starting server on port 9999")
	return a.server.ListenAndServe()
}

func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := a.server.Shutdown(ctx)
	if err != nil {
		log.Printf("Error occured while shutting down server: %v", err)
	}

	err = a.tracerShutdown(ctx)
	if err != nil {
		log.Printf("Error occured while shutting down tracer: %v", err)
	}

	return err
}
