package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// initalize database connection
func InitDB() {
	var err error
	DB, err = sql.Open("sqlite", "./prompts.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.Ping()

	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Database connection established!")

	createTables()
}

// create prompts table
func createTables() {
	promptsTable := `
	CREATE TABLE IF NOT EXISTS prompts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title TEXT NOT NULL, 
	content TEXT NOT NULL,
	category TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	deleted_at DATETIME
	);`

	_, err := DB.Exec(promptsTable)
	if err != nil {
		log.Fatal("Failed to create prmopts table:", err)
	}

	// audit table
	auditLogsTable := `
	CREATE TABLE IF NOT EXISTS audit_logs (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	admin_username TEXT NOT NULL,
	action TEXT NOT NULL,
	prompt_id INTEGER,
	prompt_title TEXT,
	timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
	details TEXT
	);`

	_, err = DB.Exec(auditLogsTable)
	if err != nil {
		log.Fatal("Failed to create audit_logs table:", err)
	}

	log.Println("Database tables created successfully!")
}
