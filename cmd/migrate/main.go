package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/nrzaman/baos-birthday-bot/internal/database"
	"github.com/nrzaman/baos-birthday-bot/util"
)

func main() {
	jsonPath := flag.String("json", "./config/birthdays.json", "Path to birthdays.json file")
	dbPath := flag.String("db", "./birthdays.db", "Path to SQLite database file")
	flag.Parse()

	fmt.Println("Birthday Bot - JSON to Database Migration Tool")
	fmt.Println("===============================================")
	fmt.Printf("JSON file: %s\n", *jsonPath)
	fmt.Printf("Database:  %s\n\n", *dbPath)

	// Read JSON file
	fmt.Println("Reading JSON file...")
	data, err := os.ReadFile(*jsonPath)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// Parse JSON
	var people util.People
	if err := json.Unmarshal(data, &people); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	fmt.Printf("Found %d birthdays in JSON file\n\n", len(people.People))

	// Open/create database
	fmt.Println("Opening database...")
	db, err := database.New(*dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()

	// Migrate data
	fmt.Println("Migrating birthdays...")
	successCount := 0
	skipCount := 0

	for _, person := range people.People {
		// Check if already exists
		existing, err := db.GetBirthday(person.Name)
		if err != nil {
			log.Printf("Warning: Error checking for existing birthday %s: %v", person.Name, err)
			continue
		}

		if existing != nil {
			fmt.Printf("  - Skipping %s (already exists)\n", person.Name)
			skipCount++
			continue
		}

		// Add birthday (with gender if present in JSON)
		if err := db.AddBirthday(person.Name, person.Birthday.Month, person.Birthday.Day, person.Gender, nil); err != nil {
			log.Printf("Warning: Failed to add birthday for %s: %v", person.Name, err)
			continue
		}

		fmt.Printf("  ✓ Migrated %s (%d/%d)\n", person.Name, person.Birthday.Month, person.Birthday.Day)
		successCount++
	}

	fmt.Printf("\nMigration complete!\n")
	fmt.Printf("  Successfully migrated: %d\n", successCount)
	fmt.Printf("  Skipped (already exist): %d\n", skipCount)
	fmt.Printf("  Total in JSON: %d\n\n", len(people.People))

	// Verify
	fmt.Println("Verifying database...")
	all, err := db.GetAllBirthdays()
	if err != nil {
		log.Fatalf("Failed to verify: %v", err)
	}

	fmt.Printf("Database now contains %d birthdays\n", len(all))
	fmt.Println("\nMigration successful! ✓")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the migrated data")
	fmt.Println("  2. Optionally add gender information using SQL:")
	fmt.Println("     UPDATE birthdays SET gender = 'male' WHERE name = 'Name';")
	fmt.Println("  3. Update main.go to use the database")
	fmt.Println("  4. Keep the JSON file as backup")
}
