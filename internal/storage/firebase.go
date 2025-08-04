package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/storage"
	"github.com/google/uuid"
	"google.golang.org/api/option"
)

type Firebase struct {
	client *storage.Client
}

func NewFirebase(projectID, credentialsPath string) *Firebase {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: projectID}
	app, err := firebase.NewApp(ctx, conf, option.WithCredentialsFile(credentialsPath))
	if err != nil {
		// Handle initialization error in production
		return &Firebase{client: nil}
	}

	client, err := app.Storage(ctx)
	if err != nil {
		return &Firebase{client: nil}
	}
	return &Firebase{client: client}
}

func (f *Firebase) Upload(ctx context.Context, filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	bucket, err := f.client.DefaultBucket()
	if err != nil {
		return "", fmt.Errorf("failed to get default bucket: %w", err)
	}

	fileName := uuid.New().String() + "-" + filePath
	obj := bucket.Object("media/" + fileName)
	w := obj.NewWriter(ctx)
	if _, err := io.Copy(w, file); err != nil {
		return "", fmt.Errorf("failed to upload to Firebase: %w", err)
	}
	if err := w.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	url, err := bucket.SignedURL("media/"+fileName, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL: %w", err)
	}
	return url, nil
}