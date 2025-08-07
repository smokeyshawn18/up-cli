package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	NeonDSN                            string
	SupabaseURL, SupabaseKey, SupabaseBucket       string
	CloudinaryCloudName, CloudinaryAPIKey, CloudinaryAPISecret string
	FirebaseProjectID, FirebaseCredentialsPath     string
}

func Load() (*Config, error) {
	// Try to load from current directory
err := godotenv.Load(".env")
	if err != nil {
		home, _ := os.UserHomeDir()
		fallbackPath := filepath.Join(home, ".config", "up-cli", ".env")
		err = godotenv.Load(fallbackPath)
		if err != nil {
			return nil, fmt.Errorf("failed to load .env file from current dir or fallback path: %w", err)
		}
	}

	cfg := &Config{
		NeonDSN:                 os.Getenv("NEON_DSN"),
		SupabaseURL:             os.Getenv("SUPABASE_URL"),
		SupabaseKey:             os.Getenv("SUPABASE_KEY"),
		SupabaseBucket:          os.Getenv("SUPABASE_BUCKET"),
		CloudinaryCloudName:     os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:        os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret:     os.Getenv("CLOUDINARY_API_SECRET"),
		FirebaseProjectID:       os.Getenv("FIREBASE_PROJECT_ID"),
		FirebaseCredentialsPath: os.Getenv("FIREBASE_CREDENTIALS_PATH"),
	}

	// Validate required variables
	if cfg.NeonDSN == "" || cfg.SupabaseURL == "" || cfg.SupabaseKey == "" || cfg.SupabaseBucket == "" ||
		cfg.CloudinaryCloudName == "" || cfg.CloudinaryAPIKey == "" || cfg.CloudinaryAPISecret == "" ||
		cfg.FirebaseProjectID == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return cfg, nil
}
