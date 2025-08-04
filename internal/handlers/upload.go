package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"up-cli/internal/database"
	"up-cli/internal/models"
	"up-cli/internal/storage"

	"github.com/google/uuid"
)

type UploadHandler struct {
	db       *database.NeonDB
	storages map[string]storage.Storage
}

func NewUploadHandler(db *database.NeonDB, storages map[string]storage.Storage) *UploadHandler {
	return &UploadHandler{db: db, storages: storages}
}

func (h *UploadHandler) Upload(ctx context.Context, filePath string) error {
	// Prompt user to select a storage provider
	fmt.Println("Select a storage provider:")
	fmt.Println("1. Supabase")
	fmt.Println("2. Cloudinary")
	fmt.Println("3. Backblaze")
	fmt.Print("Enter choice (1-3): ")

	var choice string
	_, err := fmt.Scanln(&choice)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Map choice to provider
	providerMap := map[string]string{
		"1": "supabase",
		"2": "cloudinary",
		"3": "backblaze",
	}
	provider, ok := providerMap[strings.TrimSpace(choice)]
	if !ok {
		return fmt.Errorf("invalid choice: %s, please select 1, 2, or 3", choice)
	}

	// Get storage provider
	storage, ok := h.storages[provider]
	if !ok {
		return fmt.Errorf("unsupported storage provider: %s", provider)
	}

	// Perform upload
	url, err := storage.Upload(ctx, filePath)
	if err != nil {
		return fmt.Errorf("failed to upload file to %s: %w", provider, err)
	}

	// Save metadata to database
	media := models.Media{
		ID:       uuid.New(),
		FileName: filepath.Base(filePath),
		URL:      url,
		Provider: provider,
	}

	if err := h.db.SaveMedia(media); err != nil {
		return fmt.Errorf("failed to save media metadata: %w", err)
	}

	fmt.Printf("File uploaded successfully to %s!\nURL: %s\n", provider, url)
	return nil
}