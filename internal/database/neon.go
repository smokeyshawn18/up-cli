package database

import (
	"database/sql"
	"fmt"

	"up-cli/internal/models"

	_ "github.com/lib/pq"
)

type NeonDB struct {
	db *sql.DB
}

func NewNeonDB(dsn string) (*NeonDB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Neon: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping Neon: %w", err)
	}

	// Create table if not exists
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS media (
			id UUID PRIMARY KEY,
			file_name TEXT NOT NULL,
			url TEXT NOT NULL,
			provider TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create media table: %w", err)
	}

	return &NeonDB{db: db}, nil
}

func (n *NeonDB) Close() error {
	return n.db.Close()
}

func (n *NeonDB) SaveMedia(media models.Media) error {
	_, err := n.db.Exec(
		"INSERT INTO media (id, file_name, url, provider) VALUES ($1, $2, $3, $4)",
		media.ID, media.FileName, media.URL, media.Provider,
	)
	if err != nil {
		return fmt.Errorf("failed to save media: %w", err)
	}
	return nil
}