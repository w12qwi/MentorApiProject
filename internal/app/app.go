package app

import (
	"MentorApiProject/internal/config"
	"MentorApiProject/internal/handlers/calculationHandler"
	"MentorApiProject/internal/infrastructure/db/postgres"
	"MentorApiProject/internal/service/calculationService"
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"time"
)

type App struct {
	server *http.Server
	db     *sql.DB
}

func NewApp() *App {

	postgresConfig := config.LoadPostgresConfig()

	db, err := postgres.Connect(postgresConfig)
	if err != nil {
		panic(err)
	}

	err = postgres.RunMigrations(postgresConfig)
	if err != nil {
		panic(err)
	}

	calculationRepo := postgres.NewCalculationsRepository(db)

	calculationSrvc := calculationService.NewService(calculationRepo)

	calculationHandler := calculationHandler.NewHandler(calculationSrvc)

	mux := http.NewServeMux()
	calculationHandler.RegisterRoutes(mux)

	return &App{
		server: &http.Server{
			Addr:    ":8080",
			Handler: mux,
		}, db: db,
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
	if err := a.db.Close(); err != nil {
		slog.Error("failed to close db", "error", err)
	}
	return err
}
