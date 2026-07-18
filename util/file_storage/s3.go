package file_storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Wrapper struct {
	Client *s3.Client
	Bucket string
}

func NewS3(awsRegion, awsAccessKey, awsSecretKey, awsBucket, awsEndpoint string) (*S3Wrapper, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(awsAccessKey, awsSecretKey, "")),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
		return nil, fmt.Errorf("Unable to load SDK config: %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if awsEndpoint != "" {
			o.BaseEndpoint = aws.String(awsEndpoint)
			o.UsePathStyle = true
		}
	})
	log.Println("S3 client successfully created")

	return &S3Wrapper{
		Client: client,
		Bucket: awsBucket,
	}, nil
}

func (s *S3Wrapper) UploadByte(ctx context.Context, key string, data []byte) error {
	_, err := s.Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to S3: %w", err)
	}

	return nil
}

func (s *S3Wrapper) UploadStream(ctx context.Context, key string, data io.Reader) error {
	uploader := manager.NewUploader(s.Client)

	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   data,
	})
	if err != nil {
		return fmt.Errorf("failed to upload object to S3: %w", err)
	}

	return nil
}

func (s *S3Wrapper) Download(ctx context.Context, key string) ([]byte, error) {
	resp, err := s.Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})

	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to download object from S3: %w", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	return data, nil
}

func (s *S3Wrapper) Delete(ctx context.Context, key string) error {
	_, err := s.Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object from S3: %w", err)
	}

	return nil
}
