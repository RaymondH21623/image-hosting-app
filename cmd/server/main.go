package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"shareapp/internal/application"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	err := godotenv.Load()
	if err != nil {
		logger.Error("error loading .env file", "err", err)
	}

	dataSourceName := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
	)

	conn, err := sql.Open("postgres", dataSourceName)

	if err != nil {
		logger.Error("failed to connect to db", "err", err)
		os.Exit(1)
	}

	defer conn.Close()

	// queries := db.New(conn)

	app := application.New()

	// if err := http.ListenAndServe(":8080", app.Router); err != nil {
	// 	logger.Error("server failed", "err", err)
	// 	os.Exit(1)
	// }

	err = app.Start(context.TODO())
	if err != nil {
		logger.Error("failed to start app", "err", err)
	}
}
