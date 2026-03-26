package main

import (
	"MentorApiProject/internal/app"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	godotenv.Load()

	app := app.NewApp()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	<-sigChan

	if err := app.Shutdown(); err != nil {
		log.Fatal(err)
	}

}
