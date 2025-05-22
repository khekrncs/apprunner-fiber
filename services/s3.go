// services/s3.go
package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"

	appConfig "github.com/khekrn/apprunner-fiber/config"
	"github.com/khekrn/apprunner-fiber/models"
)

type S3Service struct {
	client *s3.Client
	config *appConfig.S3Config
	mu     sync.RWMutex
	once   sync.Once
}

func NewS3Service() *S3Service {
	return &S3Service{
		config: appConfig.GetS3Config(),
	}
}

// Lazy initialization of S3 client
func (s *S3Service) getClient() (*s3.Client, error) {
	s.mu.RLock()
	if s.client != nil {
		defer s.mu.RUnlock()
		return s.client, nil
	}
	s.mu.RUnlock()

	var err error
	s.once.Do(func() {
		s.mu.Lock()
		defer s.mu.Unlock()

		var cfg aws.Config

		// Build configuration options
		configOptions := []func(*config.LoadOptions) error{
			config.WithRegion(s.config.Region),
		}

		// Add retry configuration
		if s.config.MaxRetries > 0 {
			configOptions = append(configOptions, config.WithRetryMaxAttempts(s.config.MaxRetries))
		}

		// Load default config using AWS credential provider chain
		cfg, err = config.LoadDefaultConfig(context.TODO(), configOptions...)
		if err != nil {
			return
		}

		s.client = s3.NewFromConfig(cfg)
	})

	return s.client, err
}

func (s *S3Service) PutObject(ctx context.Context, key string, data []byte, contentType string, metadata map[string]string) (*models.FileInfo, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get S3 client: %w", err)
	}

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.config.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
		Metadata:    metadata,
	}

	result, err := client.PutObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to put object: %w", err)
	}

	fileInfo := &models.FileInfo{
		Key:          key,
		ContentType:  contentType,
		Size:         int64(len(data)),
		ETag:         aws.ToString(result.ETag),
		Metadata:     metadata,
		UploadedAt:   time.Now(),
		LastModified: time.Now(),
	}

	return fileInfo, nil
}

func (s *S3Service) GetObject(ctx context.Context, key string) ([]byte, *models.FileInfo, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get S3 client: %w", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	}

	result, err := client.GetObject(ctx, input)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read object body: %w", err)
	}

	fileInfo := &models.FileInfo{
		Key:          key,
		ContentType:  aws.ToString(result.ContentType),
		Size:         aws.ToInt64(result.ContentLength),
		ETag:         aws.ToString(result.ETag),
		Metadata:     result.Metadata,
		LastModified: aws.ToTime(result.LastModified),
	}

	return data, fileInfo, nil
}

func (s *S3Service) DeleteObject(ctx context.Context, key string) error {
	client, err := s.getClient()
	if err != nil {
		return fmt.Errorf("failed to get S3 client: %w", err)
	}

	input := &s3.DeleteObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	}

	_, err = client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func (s *S3Service) ListObjects(ctx context.Context, prefix string) ([]*models.FileInfo, error) {
	client, err := s.getClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get S3 client: %w", err)
	}

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.config.Bucket),
		Prefix: aws.String(prefix),
	}

	result, err := client.ListObjectsV2(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %w", err)
	}

	var files []*models.FileInfo
	for _, obj := range result.Contents {
		fileInfo := &models.FileInfo{
			Key:          aws.ToString(obj.Key),
			Size:         aws.ToInt64(obj.Size),
			ETag:         aws.ToString(obj.ETag),
			LastModified: aws.ToTime(obj.LastModified),
		}
		files = append(files, fileInfo)
	}

	return files, nil
}

func (s *S3Service) GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	if !s.config.EnablePresignedURLs {
		return "", fmt.Errorf("presigned URLs are disabled")
	}

	client, err := s.getClient()
	if err != nil {
		return "", fmt.Errorf("failed to get S3 client: %w", err)
	}

	presignClient := s3.NewPresignClient(client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

func (s *S3Service) GeneratePresignedUploadURL(ctx context.Context, key string, expiration time.Duration) (string, error) {
	if !s.config.EnablePresignedURLs {
		return "", fmt.Errorf("presigned URLs are disabled")
	}

	client, err := s.getClient()
	if err != nil {
		return "", fmt.Errorf("failed to get S3 client: %w", err)
	}

	presignClient := s3.NewPresignClient(client)

	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiration
	})

	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return request.URL, nil
}

func (s *S3Service) ObjectExists(ctx context.Context, key string) (bool, error) {
	client, err := s.getClient()
	if err != nil {
		return false, fmt.Errorf("failed to get S3 client: %w", err)
	}

	input := &s3.HeadObjectInput{
		Bucket: aws.String(s.config.Bucket),
		Key:    aws.String(key),
	}
	_, err = client.HeadObject(ctx, input)
	if err != nil {
		var noSuchKey *types.NoSuchKey
		if errors.As(err, &noSuchKey) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check object existence: %w", err)
	}

	return true, nil
}
