package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"up-cli/internal/config"
	"up-cli/internal/database"
	"up-cli/internal/handlers"
	"up-cli/internal/storage"

	"github.com/spf13/cobra"
)

var version = "dev" // This will be overridden at build time using -ldflags

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := database.NewNeonDB(cfg.NeonDSN)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	supabaseStorage, err := storage.NewSupabaseFromEnv()
	if err != nil {
		log.Fatalf("Failed to initialize Supabase storage: %v", err)
	}

	cloudinaryStorage := storage.NewCloudinary(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)

	backblazeStorage, err := storage.NewBackblazeFromEnv(context.Background())
	if err != nil {
		log.Fatalf("Failed to initialize Backblaze storage: %v", err)
	}

	storages := map[string]storage.Storage{
		"supabase":   supabaseStorage,
		"cloudinary": cloudinaryStorage,
		"backblaze":  backblazeStorage,
	}

	handler := handlers.NewUploadHandler(db, storages)

	rootCmd := &cobra.Command{
		Use:     "up-cli",
		Short:   "Upload media to cloud storage",
		Version: version,
	}

uploadCmd := &cobra.Command{
    Use:   "upload [files...]",
    Short: "Upload one or more files to a cloud storage provider",
    Args:  cobra.MinimumNArgs(1), // At least one file required
    RunE: func(cmd *cobra.Command, args []string) error {
        return handler.Upload(cmd.Context(), args) // pass []string of file paths
    },
}

	rootCmd.AddCommand(uploadCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Command error: %v\n", err)
		os.Exit(1)
	}
}
