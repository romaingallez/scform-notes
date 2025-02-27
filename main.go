package main

import (
	"encoding/json"
	"log"
	"os"
)

func init() {
	// setup logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	scformURL := "https://mewo.sc-form.net/"
	username := "GALLEZ"
	password := "[REDACTED]"

	student, err := GetStudentGrades(scformURL, username, password)
	if err != nil {
		log.Fatal("Failed to get student grades:", err)
	}

	// Create the JSON output
	jsonData, err := json.MarshalIndent(student, "", "  ")
	if err != nil {
		log.Fatal("Failed to marshal student data to JSON:", err)
	}

	err = os.WriteFile("grades.json", jsonData, 0644)
	if err != nil {
		log.Fatal("Failed to write grades to file:", err)
	}
}
