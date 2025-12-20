package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"shareapp/utils"
	"time"

	_ "github.com/lib/pq"

	"shareapp/internal/data"

	"github.com/joho/godotenv"
)

const version = "1.0.0"

type config struct {
	port int    `env:"SERVER_PORT"`
	env  string `env:"ENVIRONMENT"`
	db   struct {
		dsn string `env:"POSTGRES_URL"`
	}
}

type application struct {
	db       *sql.DB
	config   config
	queries  *data.Queries
	jwtMaker *utils.JWTMaker
	logger   *slog.Logger
	// minio    *minio.Client
}

func main() {

	var cfg config

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.StringVar(&cfg.db.dsn, "dsn", os.Getenv("POSTGRES_URL"), "PostgreSQL DSN")
	flag.Parse()

	// if err := env.Parse(&cfg); err != nil {
	// 	log.Fatalf("Failed to read the environment variables: %v", err)
	// 	return
	// }

	// cfg.port = os.Getenv("SERVER_PORT")

	// dsn := os.Getenv("POSTGRES_URL")

	// port := os.Getenv("SERVER_PORT")

	ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	// minioClient, err := minio.New("localhost:9000", &minio.Options{
	// 	Creds:  credentials.NewStaticV4("minio", "minio123", ""),
	// 	Secure: false,
	// })

	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	bucketName := "media"
	location := "us-east-1"

	// err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	// if err != nil {
	// 	// Check to see if we already own this bucket (which happens if you run this twice)
	// 	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	// 	if errBucketExists == nil && exists {
	// 		log.Printf("We already own %s\n", bucketName)
	// 	} else {
	// 		log.Fatalln(err)
	// 	}
	// } else {
	// 	log.Printf("Successfully created %s\n", bucketName)
	// }

	app := &application{
		config:   cfg,
		db:       db,
		queries:  data.New(db),
		jwtMaker: utils.NewJWTMaker("secret-key"),
		// minio:    minioClient,
		logger: logger,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	fmt.Println(cfg.port)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	return db, nil
}
