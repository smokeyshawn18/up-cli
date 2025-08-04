package storage

import "context"

type Storage interface {
	Upload(ctx context.Context, filePath string) (string, error)
}