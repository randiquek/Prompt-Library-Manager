package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"prompt-library/database"
	"prompt-library/models"

	"github.com/gorilla/mux"
)

// function to convert string ID to int
func parseID(id string) int {
	var idInt int
	fmt.Sscanf(id, "%d", &idInt)
	return idInt
}

// Gets all prompts from database
func GetAllPrompts(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
	SELECT id, title, content, category, created_at, updated_at
	FROM prompts
	WHERE deleted_at IS NULL
	ORDER BY created_at DESC
	`)
	if err != nil {
		log.Println("Error querying prompts:", err)
		http.Error(w, "Failed to fetch prompts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// empty slice for holding all prompts
	prompts := []models.Prompt{}

	// Loops through each row and check each prompt struct

	for rows.Next() {
		var p models.Prompt
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Category, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			log.Println("Error scanning prompt:", err)
			continue
		}
		prompts = append(prompts, p)
	}

	// Set response header to JSON
	w.Header().Set("Content-Type", "application/json")

	// convert prompts to JSON and send response
	json.NewEncoder(w).Encode(prompts)
}

func CreatePrompt(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body

	var input models.PromptInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields for PromptInput
	if input.Title == "" || input.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Insert parsed data into database
	result, err := database.DB.Exec(`
	INSERT INTO prompts (title, content, category)
	VALUES (?, ?, ?)
	`, input.Title, input.Content, input.Category)

	if err != nil {
		log.Println("Error creating prompt:", err)
		http.Error(w, "Failed to create prompt", http.StatusInternalServerError)
		return
	}

	// Get the ID of new prompt
	id, err := result.LastInsertId()
	if err != nil {
		log.Println("Error getting last ID insert:", err)
		http.Error(w, "Failed to get prompt ID", http.StatusInternalServerError)
		return
	}

	// get username from context
	username := r.Context().Value("username").(string)

	// TODO: Add audit log
	database.LogAudit(username, "CREATE", int(id), input.Title, "")

	// Return the created prompt with new ID
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201 status
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":      id,
		"message": "Prompt created successfully",
	})

}

func UpdatePrompt(w http.ResponseWriter, r *http.Request) {
	// get ID from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Parse JSON request body

	var input models.PromptInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if input.Title == "" || input.Content == "" {
		http.Error(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	// Update database
	result, err := database.DB.Exec(`
	UPDATE prompts
	SET title = ?, content = ?, category = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`, input.Title, input.Content, input.Category, id)

	if err != nil {
		log.Println("Error updating prompt:", err)
		http.Error(w, "Failed to update prompt", http.StatusInternalServerError)
		return
	}

	// check if any rows were affected, like if prompt already existed
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected", err)
		http.Error(w, "Failed to update prompt", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Prompt not found", http.StatusNotFound)
		return
	}

	// get username from context
	username := r.Context().Value("username").(string)

	// TODO: Add audit log
	database.LogAudit(username, "UPDATE", parseID(id), input.Title, "")

	// Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Prompt updated successfully",
	})

}

func DeletePrompt(w http.ResponseWriter, r *http.Request) {
	// get id from URL
	vars := mux.Vars(r)
	id := vars["id"]

	// Get prompt details for audit logging
	var title string
	err := database.DB.QueryRow("SELECT title FROM prompts WHERE id = ?", id).Scan(&title)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Prompt not found", http.StatusNotFound)
			return
		}
		log.Println("Error fetching prompt", err)
		http.Error(w, "Failed to delete prompt", http.StatusInternalServerError)
		return
	}

	// Delete from Database (Soft delete)
	result, err := database.DB.Exec(`
	UPDATE prompts 
	SET deleted_at = CURRENT_TIMESTAMP 
	WHERE id = ? AND deleted_at IS NULL`, id)
	if err != nil {
		log.Println("Error deleting prompt:", err)
		http.Error(w, "Failed to delete prompt", http.StatusInternalServerError)
		return
	}

	// Check if any affected rows
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error checking rows affected:", err)
		http.Error(w, "Failed to delete prompt", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Prompt not found", http.StatusNotFound)
		return
	}
	// get username from context
	username := r.Context().Value("username").(string)

	// TODO: Add audit log
	database.LogAudit(username, "SOFT_DELETE", parseID(id), title, "Marked as deleted")

	// Return success
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Prompt deleted successfully!",
	})

}

// returns audit log entries
func GetAuditLogs(w http.ResponseWriter, r *http.Request) {
	rows, err := database.DB.Query(`
	SELECT id, admin_username, action, prompt_id, prompt_title, timestamp, details
	FROM audit_logs
	ORDER BY timestamp DESC
	`)
	if err != nil {
		log.Println("Error querying audit logs", err)
		http.Error(w, "Failed to fetch audit logs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	logs := []models.AuditLog{}

	for rows.Next() {
		var auditLog models.AuditLog
		err := rows.Scan(&auditLog.ID, &auditLog.AdminUsername, &auditLog.Action, &auditLog.PromptID, &auditLog.PromptTitle, &auditLog.Timestamp, &auditLog.Details)
		if err != nil {
			log.Println("Error scanning audit log:", err)
			continue
		}
		logs = append(logs, auditLog)

	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
	
}