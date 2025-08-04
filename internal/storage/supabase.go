package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Supabase struct {
	projectURL string
	apiKey     string
	bucket     string
}

// Constructor loading config from .env
func NewSupabaseFromEnv() (*Supabase, error) {
	_ = godotenv.Load()

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	bucket := os.Getenv("SUPABASE_BUCKET")

	if url == "" || key == "" || bucket == "" {
		return nil, fmt.Errorf("missing one or more Supabase env variables")
	}

	return &Supabase{
		projectURL: url,
		apiKey:     key,
		bucket:     bucket,
	}, nil
}

// Upload implements uploading file to Supabase Storage bucket
func (s *Supabase) Upload(ctx context.Context, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}

	fileName := filepath.Base(filePath)

	// Detect MIME type from file extension
	ext := filepath.Ext(fileName)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream" // fallback
	}

	uploadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", s.projectURL, s.bucket, fileName)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uploadURL, bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+s.apiKey)
	req.Header.Set("Content-Type", mimeType)
	req.Header.Set("x-upsert", "true")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("upload error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", body)
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", s.projectURL, s.bucket, fileName)
	return publicURL, nil
}
