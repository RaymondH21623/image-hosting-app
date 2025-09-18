package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"shareapp/internal/server"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := os.Getenv("POSTGRES_URL")

	port := os.Getenv("SERVER_PORT")

	srv, err := server.New(port, dsn)
	if err != nil {
		log.Fatalf("failed to start server: %v", err)
	}

	fmt.Printf("Starting server at port %s\n", port)
	log.Fatal(http.ListenAndServe(port, srv))
}
