package services

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Service struct {
	bucket string
	region string
	client *s3.Client
}

func NewS3Service(ctx context.Context, bucket, region string) (*S3Service, error) {
	bucket = strings.TrimSpace(bucket)
	region = strings.TrimSpace(region)
	if bucket == "" {
		return nil, fmt.Errorf("S3_BUCKET_NAME is required")
	}
	if region == "" {
		return nil, fmt.Errorf("AWS_REGION is required")
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &S3Service{
		bucket: bucket,
		region: region,
		client: s3.NewFromConfig(cfg),
	}, nil
}

// UploadPublic uploads bytes to S3 and returns (key, publicURL).
// Bucket must allow public reads for the returned URL to be reachable.
func (s *S3Service) UploadPublic(ctx context.Context, key string, body []byte, contentType string) (string, string, error) {
	key = strings.TrimPrefix(strings.TrimSpace(key), "/")
	if key == "" {
		return "", "", fmt.Errorf("key is required")
	}

	input := &s3.PutObjectInput{
		Bucket:       aws.String(s.bucket),
		Key:          aws.String(key),
		Body:         bytes.NewReader(body),
		ContentType:  aws.String(contentType),
		CacheControl: aws.String("public, max-age=31536000, immutable"),
	}

	if _, err := s.client.PutObject(ctx, input); err != nil {
		return "", "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Standard virtual-hosted–style URL for us-east-1.
	u := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.s3.%s.amazonaws.com", s.bucket, s.region),
		Path:   path.Clean("/" + key),
	}
	return key, u.String(), nil
}

// Delete removes an object by key.
func (s *S3Service) Delete(ctx context.Context, key string) error {
	key = strings.TrimPrefix(strings.TrimSpace(key), "/")
	if key == "" {
		return nil
	}
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}

// ReadAllWithLimit reads from r up to maxBytes (inclusive) and errors if exceeded.
func ReadAllWithLimit(r io.Reader, maxBytes int64) ([]byte, error) {
	lr := &io.LimitedReader{R: r, N: maxBytes + 1}
	b, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}
	if int64(len(b)) > maxBytes {
		return nil, fmt.Errorf("file too large (max %d bytes)", maxBytes)
	}
	return b, nil
}
