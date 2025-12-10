package models

import (
	"time"
)

// struct for admin log records

type AuditLog struct {
	ID int `json:"id"`
	AdminUsername string `json:"admin_username"`
	Action string `json:"action"`
	PromptID int `json:"prompt_id"`
	PromptTitle string `json:"prompt_title"`
	Timestamp time.Time `json:"timestamp"`
	Details string `json:"details"`
}