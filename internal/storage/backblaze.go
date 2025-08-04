package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/kurin/blazer/b2"
)

type Backblaze struct {
	client *b2.Client
	bucket *b2.Bucket
}

func NewBackblaze(ctx context.Context, keyID, appKey, bucketName string) (*Backblaze, error) {
	client, err := b2.NewClient(ctx, keyID, appKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create b2 client: %w", err)
	}

	bucket, err := client.Bucket(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to get bucket: %w", err)
	}

	return &Backblaze{
		client: client,
		bucket: bucket,
	}, nil
}

func NewBackblazeFromEnv(ctx context.Context) (*Backblaze, error) {
	keyID := os.Getenv("B2_KEY_ID")
	appKey := os.Getenv("B2_APPLICATION_KEY")
	bucket := os.Getenv("B2_BUCKET")

	if keyID == "" || appKey == "" || bucket == "" {
		return nil, fmt.Errorf("missing one or more Backblaze env variables")
	}

	return NewBackblaze(ctx, keyID, appKey, bucket)
}

func (b *Backblaze) Upload(ctx context.Context, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	fileName := filepath.Base(filePath)

	writer := b.bucket.Object(fileName).NewWriter(ctx)
	if _, err := io.Copy(writer, file); err != nil {
		writer.Close()
		return "", fmt.Errorf("upload copy failed: %w", err)
	}

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("upload close failed: %w", err)
	}

	// Construct public URL (adjust if your bucket is public)
	publicURL := fmt.Sprintf("https://f000.backblazeb2.com/file/%s/%s", b.bucket.Name(), fileName)
	return publicURL, nil
}
