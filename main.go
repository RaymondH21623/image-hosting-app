package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"shareapp/internal/db"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
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

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	conn, err := sql.Open("postgres", dataSourceName)
	if err := conn.PingContext(ctx); err != nil {
		logger.Error("failed to connect to db", "err", err)
		os.Exit(1)
	}
	defer conn.Close()

	queries := db.New(conn)

	user, err := queries.CreateUser(ctx, db.CreateUserParams{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpw",
	})
	if err != nil {
		logger.Error("failed to create user", "err", err)
		os.Exit(1)
	}

	logger.Info("created user", "id", user.ID, "email", user.Email)
}
