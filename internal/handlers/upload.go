package handlers

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"

	"up-cli/internal/database"
	"up-cli/internal/models"
	"up-cli/internal/storage"
)

// UploadHandler is responsible for handling file uploads to multiple cloud storage providers.
// It holds a reference to the database for saving metadata and a map of storage clients.
type UploadHandler struct {
    db       *database.NeonDB
    storages map[string]storage.Storage
}

// NewUploadHandler constructs a new UploadHandler with a connected database and available storages.
func NewUploadHandler(db *database.NeonDB, storages map[string]storage.Storage) *UploadHandler {
    return &UploadHandler{db: db, storages: storages}
}

// Upload uploads multiple files concurrently to the selected storage provider.
// It prompts the user to select the provider, shows a progress bar,
// uploads files with a concurrency limit, saves metadata, and collects errors.
func (h *UploadHandler) Upload(ctx context.Context, filePaths []string) error {
    if len(filePaths) == 0 {
        return fmt.Errorf("no files provided for upload")
    }

    // Prompt user to select a storage provider only once.
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

    // Map user input to provider name keys.
    providerMap := map[string]string{
        "1": "supabase",
        "2": "cloudinary",
        "3": "backblaze",
    }
    provider, ok := providerMap[strings.TrimSpace(choice)]
    if !ok {
        return fmt.Errorf("invalid choice: %s, please select 1, 2, or 3", choice)
    }

    storage, ok := h.storages[provider]
    if !ok {
        return fmt.Errorf("unsupported storage provider: %s", provider)
    }

    // Initialize a progress bar to show overall upload progress.
    bar := progressbar.Default(int64(len(filePaths)))

    // Mutex to protect error slice when accessed from multiple goroutines.
    var errsMu sync.Mutex
    var errs []string
    var wg sync.WaitGroup

    // Set concurrency limit to avoid overwhelming resources or API limits.
    concurrencyLimit := 5
    sem := make(chan struct{}, concurrencyLimit)

    // Iterate over all file paths and upload concurrently.
    for _, filePath := range filePaths {
        wg.Add(1)
        sem <- struct{}{} // Acquire a slot for concurrency limiting.

        go func(fp string) {
            defer wg.Done()
            defer func() { <-sem }() // Release the slot.

            // Upload file to the selected storage provider.
            url, err := storage.Upload(ctx, fp)
            if err != nil {
                errsMu.Lock()
                errs = append(errs, fmt.Sprintf("failed to upload %s: %v", fp, err))
                errsMu.Unlock()
                bar.Add(1) // Still advance progress bar.
                return
            }

            // Create media metadata for the uploaded file.
            media := models.Media{
                ID:       uuid.New(),
                FileName: filepath.Base(fp),
                URL:      url,
                Provider: provider,
            }

            // Save media metadata in the database.
            if err := h.db.SaveMedia(media); err != nil {
                errsMu.Lock()
                errs = append(errs, fmt.Sprintf("failed to save metadata for %s: %v", fp, err))
                errsMu.Unlock()
                bar.Add(1)
                return
            }

            // Advance progress bar and print success message.
            bar.Add(1)
            fmt.Printf("\nFile uploaded successfully to %s!\nFile: %s\nURL: %s\n\n", provider, fp, url)
        }(filePath)
    }

    // Wait for all uploads to complete.
    wg.Wait()

    // If any errors occurred during uploads or saving metadata, return them.
    if len(errs) > 0 {
        return fmt.Errorf("some errors occurred:\n%s", strings.Join(errs, "\n"))
    }

    return nil
}
