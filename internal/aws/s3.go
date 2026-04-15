package storage

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Endpoint        string
	UsePathStyle    bool
}

type S3Storage struct {
	Client        *s3.Client
	PresignClient *s3.PresignClient
	Bucket        string
}

func New(ctx context.Context, cfg Config, bucket string) (*S3Storage, error) {
	awsCFG, err := loadAWSConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(awsCFG, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(cfg.Endpoint)
		o.UsePathStyle = cfg.UsePathStyle
	})

	presignClient := s3.NewPresignClient(s3Client)

	storage := &S3Storage{
		Client:        s3Client,
		PresignClient: presignClient,
		Bucket:        bucket,
	}

	// if err := storage.verifyCredentials(ctx, awsCFG); err != nil {
	// 	return nil, err
	// }

	if err := storage.ensureBucket(ctx); err != nil {
		return nil, err
	}

	return storage, nil
}

func loadAWSConfig(ctx context.Context, cfg Config) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKeyID,
				cfg.SecretAccessKey,
				"",
			),
		),
	)
}

func (s *S3Storage) verifyCredentials(ctx context.Context, cfg aws.Config) (string, error) {
	creds, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		return "", err
	}

	return creds.AccessKeyID, nil
}

func (s *S3Storage) ensureBucket(ctx context.Context) error {
	createCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_, err := s.Client.CreateBucket(createCtx, &s3.CreateBucketInput{
		Bucket: aws.String(s.Bucket),
	})

	return err
}
