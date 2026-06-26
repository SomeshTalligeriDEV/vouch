package external

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// R2Presigner implements service.Presigner against Cloudflare R2 using its
// S3-compatible API.
type R2Presigner struct {
	client    *s3.Client
	presign   *s3.PresignClient
	bucket    string
	publicURL string
}

// NewR2Presigner constructs an R2Presigner. publicURL is the base URL objects
// are served from (e.g. a bucket public domain or custom CDN domain); if empty
// it falls back to the endpoint/bucket path.
func NewR2Presigner(accessKey, secretKey, endpoint, bucket, publicURL string) *R2Presigner {
	client := s3.New(s3.Options{
		Region:       "auto",
		BaseEndpoint: aws.String(endpoint),
		Credentials:  credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		UsePathStyle: true,
	})
	if publicURL == "" {
		publicURL = strings.TrimRight(endpoint, "/") + "/" + bucket
	}
	return &R2Presigner{
		client:    client,
		presign:   s3.NewPresignClient(client),
		bucket:    bucket,
		publicURL: strings.TrimRight(publicURL, "/"),
	}
}

// PresignPut returns a presigned PUT URL and the resulting public URL.
func (r *R2Presigner) PresignPut(ctx context.Context, key, contentType string, ttl time.Duration) (string, string, error) {
	req, err := r.presign.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}, s3.WithPresignExpires(ttl))
	if err != nil {
		return "", "", fmt.Errorf("R2Presigner.PresignPut: %w", err)
	}
	publicURL := r.publicURL + "/" + key
	return req.URL, publicURL, nil
}
