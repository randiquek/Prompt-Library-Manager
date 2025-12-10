package models

import (
	"time"
)

// Prompt struct (class or blueprint)
type Prompt struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Content string `json:"content"`
	Category string `json:"category"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Struct used when creating/updating prompts (no ID or timestamps)
type PromptInput struct {
	Title string `json:"title"`
	Content string `json:"content"`
	Category string `json:"category"`
}