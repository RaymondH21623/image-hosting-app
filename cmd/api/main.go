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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	_ "github.com/lib/pq"

	"shareapp/internal/data"

	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Config struct {
	port int    `env:"SERVER_PORT"`
	env  string `env:"ENVIRONMENT"`
	db   struct {
		dsn string `env:"POSTGRES_URL"`
	}
}

type application struct {
	db            *sql.DB
	config        Config
	queries       *data.Queries
	jwtMaker      *utils.JWTMaker
	logger        *slog.Logger
	S3Client      *s3.Client
	presignClient *s3.PresignClient
}

func main() {

	var cfg Config

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

	//ctx := context.Background()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	S3Config, err := loadAWSConfig(context.TODO())
	if err != nil {
		logger.Error("unable to load AWS SDK config, " + err.Error())
		os.Exit(1)
	}

	s3Client := s3.NewFromConfig(S3Config, func(o *s3.Options) {
		o.BaseEndpoint = aws.String("http://localhost:3900")
		o.UsePathStyle = true

	})

	presigner := s3.NewPresignClient(s3Client)

	ctx := context.Background()

	creds, err := S3Config.Credentials.Retrieve(ctx)
	if err != nil {
		logger.Error("unable to retrieve AWS credentials, " + err.Error())
		os.Exit(1)
	}

	logger.Info("aws access key id: " + creds.AccessKeyID)

	createCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err = s3Client.CreateBucket(createCtx, &s3.CreateBucketInput{
		Bucket: aws.String("media"),
	})

	if err != nil {
		logger.Error("unable to create bucket, " + err.Error())
		os.Exit(1)
	}

	logger.Info("bucket ready", "bucket", "media")

	app := &application{
		config:        cfg,
		db:            db,
		queries:       data.New(db),
		jwtMaker:      utils.NewJWTMaker("secret-key"),
		S3Client:      s3Client,
		presignClient: presigner,
		logger:        logger,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(cfg Config) (*sql.DB, error) {
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

func loadAWSConfig(ctx context.Context) (aws.Config, error) {
	return config.LoadDefaultConfig(
		ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"),
				"",
			),
		),
	)
}
