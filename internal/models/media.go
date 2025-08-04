package models

import "github.com/google/uuid"

type Media struct {
	ID        uuid.UUID `json:"id"`
	FileName  string    `json:"file_name"`
	URL       string    `json:"url"`
	Provider  string    `json:"provider"`
}