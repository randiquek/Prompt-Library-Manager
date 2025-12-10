// TEMP FILE - used for test data, delete before production

package main

import (
	"log"
	"prompt-library/database"
)

func main() {
	database.InitDB()

	// test prompt
	_, err := database.DB.Exec(`
		INSERT INTO prompts (title, content, category)
		VALUES
			('Write a blog post', 'Write an engaging blog post about [TOPIC]. Include an introduction, 3 main points, and a conclusion.', 'Writing'),
			('Code review', 'Review this code and suggest improvements: [CODE]', 'Coding'),
			('Email draft', 'Draft a professional email to [RECIPIENT] about [SUBJECT]', 'Business') 
	`)

	if err != nil {
		log.Fatal("Error inserting test data:", err)
	}
	log.Println("Test data inserted successfully!")
}