package database

import (
	"log"
)

 // Records an admin action to the audit_logs table

func LogAudit(username, action string, promptID int, promptTitle, details string) error {
	_, err := DB.Exec(`
	INSERT INTO audit_logs (admin_username, action, prompt_id, prompt_title, details)
	VALUES (?, ?, ?, ?, ?)
	`, username, action, promptID, promptTitle, details)

	if err != nil {
		log.Printf("Error logging audit: %v", err)
		return err
	}

	log.Printf("Audit logged: %s by %s on prompt %d", action, username, promptID)
	return nil
}