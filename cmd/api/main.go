package main

import (
	"context"
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"os"
	storage "shareapp/internal/aws"
	"shareapp/internal/data"
	"shareapp/internal/mailer"
	"shareapp/utils"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

const version = "1.0.0"

type Config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
}

type application struct {
	config    Config
	models    data.Models
	jwtMaker  *utils.JWTMaker
	logger    *slog.Logger
	S3Storage *storage.S3Storage
	mailer    mailer.Mailer
	wg        sync.WaitGroup
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

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.smtp.host, "smtp-host", "sandbox.smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Echo <no-reply@echo.raymondhoang.net>", "SMTP sender")

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool established")

	// S3Config, err := loadAWSConfig(context.TODO())
	// if err != nil {
	// 	logger.Error("unable to load AWS SDK config, " + err.Error())
	// 	os.Exit(1)
	// }

	// s3Client := s3.NewFromConfig(S3Config, func(o *s3.Options) {
	// 	o.BaseEndpoint = aws.String("http://localhost:3900")
	// 	o.UsePathStyle = true

	// })

	// presigner := s3.NewPresignClient(s3Client)

	// ctx := context.Background()

	// creds, err := S3Config.Credentials.Retrieve(ctx)
	// if err != nil {
	// 	logger.Error("unable to retrieve AWS credentials, " + err.Error())
	// 	os.Exit(1)
	// }

	// logger.Info("aws access key id: " + creds.AccessKeyID)

	// createCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	// defer cancel()

	// _, err = s3Client.CreateBucket(createCtx, &s3.CreateBucketInput{
	// 	Bucket: aws.String("media"),
	// })

	// if err != nil {
	// 	logger.Error("unable to create bucket, " + err.Error())
	// 	os.Exit(1)
	// }

	// logger.Info("bucket ready", "bucket", "media")

	ctx := context.Background()

	s3Config := storage.Config{
		AccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		SecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		Endpoint:        "http://localhost:3900",
		Region:          "us-east-1",
		UsePathStyle:    true,
	}

	storage, err := storage.New(ctx, s3Config, "media")
	if err != nil {
		logger.Error("unable to create S3 client, " + err.Error())
		os.Exit(1)
	}

	app := &application{
		config:    cfg,
		models:    data.NewModels(db),
		jwtMaker:  utils.NewJWTMaker("secret-key"),
		S3Storage: storage,
		logger:    logger,
		mailer:    mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

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
