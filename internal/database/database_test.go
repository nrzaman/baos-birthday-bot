package database_test

import (
	"testing"

	"github.com/nrzaman/baos-birthday-bot/internal/database"
)

// Helper function to create an in-memory test database
func setupTestDB(t *testing.T) *database.DB {
	t.Helper()
	// Use :memory: for in-memory database (perfect for tests!)
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Close() // Best effort close in tests
	})
	return db
}

func TestAddAndGetBirthday(t *testing.T) {
	db := setupTestDB(t)
	expected_month := 1
	expected_day := 25

	// Add a birthday
	gender := "female"
	err := db.AddBirthday("Alice", expected_month, expected_day, &gender, nil)
	if err != nil {
		t.Fatalf("Failed to add birthday: %v", err)
	}

	// Get the birthday
	birthday, err := db.GetBirthday("Alice")
	if err != nil {
		t.Fatalf("Failed to get birthday: %v", err)
	}

	if birthday == nil {
		t.Fatal("Birthday not found")
	}

	if birthday.Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", birthday.Name)
	}

	if birthday.Month != expected_month {
		t.Errorf("Expected month %d, got %d", expected_month, birthday.Month)
	}

	if birthday.Day != expected_day {
		t.Errorf("Expected day %d, got %d", expected_day, birthday.Day)
	}

	if birthday.Gender == nil || *birthday.Gender != "female" {
		t.Error("Expected gender 'female'")
	}
}

func TestGetBirthdaysByMonth(t *testing.T) {
	db := setupTestDB(t)
	expected_month_alice_ben := 1
	expected_day_alice := 25
	expected_day_ben := 31
	expected_month_cassidy := 12
	expected_day_cassidy := 2

	// Add multiple birthdays
	_ = db.AddBirthday("Alice", expected_month_alice_ben, expected_day_alice, nil, nil)
	_ = db.AddBirthday("Bob", expected_month_alice_ben, expected_day_ben, nil, nil)
	_ = db.AddBirthday("Cassidy", expected_month_cassidy, expected_day_cassidy, nil, nil)

	// Get March birthdays
	birthdays, err := db.GetBirthdaysByMonth(expected_month_alice_ben)
	if err != nil {
		t.Fatalf("Failed to get birthdays: %v", err)
	}

	if len(birthdays) != 2 {
		t.Errorf("Expected 2 birthdays in January, got %d", len(birthdays))
	}
}

func TestGetBirthdaysByDate(t *testing.T) {
	db := setupTestDB(t)
	expected_month_alice := 1
	expected_day_alice := 25
	expected_month_bruce_cassidy := 12
	expected_day_bruce_cassidy := 2

	// Add birthdays
	_ = db.AddBirthday("Alice", expected_month_alice, expected_day_alice, nil, nil)
	_ = db.AddBirthday("Bruce", expected_month_bruce_cassidy, expected_day_bruce_cassidy, nil, nil) // Same day!
	_ = db.AddBirthday("Cassidy", expected_month_bruce_cassidy, expected_day_bruce_cassidy, nil, nil)

	// Get birthdays on March 15
	birthdays, err := db.GetBirthdaysByDate(expected_month_bruce_cassidy, expected_day_bruce_cassidy)
	if err != nil {
		t.Fatalf("Failed to get birthdays: %v", err)
	}

	if len(birthdays) != 2 {
		t.Errorf("Expected 2 birthdays on December 2, got %d", len(birthdays))
	}
}

func TestUpdateBirthday(t *testing.T) {
	db := setupTestDB(t)
	expected_month_alice := 1
	expected_day_alice := 25

	// Add a birthday
	_ = db.AddBirthday("Alice", expected_month_alice, expected_day_alice, nil, nil)

	// Update it
	gender := "female"
	err := db.UpdateBirthday("Alice", expected_month_alice, (expected_day_alice + 1), &gender, nil)
	if err != nil {
		t.Fatalf("Failed to update birthday: %v", err)
	}

	// Verify
	birthday, _ := db.GetBirthday("Alice")
	if birthday.Day != (expected_day_alice + 1) {
		t.Errorf("Expected day %d, got %d", (expected_day_alice + 1), birthday.Day)
	}
	if birthday.Gender == nil || *birthday.Gender != "female" {
		t.Error("Expected gender 'female'")
	}
}

func TestDeleteBirthday(t *testing.T) {
	db := setupTestDB(t)

	// Add a birthday
	_ = db.AddBirthday("Alice", 1, 25, nil, nil)

	// Delete it
	err := db.DeleteBirthday("Alice")
	if err != nil {
		t.Fatalf("Failed to delete birthday: %v", err)
	}

	// Verify it's gone
	birthday, _ := db.GetBirthday("Alice")
	if birthday != nil {
		t.Error("Birthday should have been deleted")
	}
}

func TestGetPronoun(t *testing.T) {
	tests := []struct {
		name        string
		gender      *string
		subjectForm bool
		expected    string
	}{
		{"Male subject", stringPtr("male"), true, "he"},
		{"Male possessive", stringPtr("male"), false, "his"},
		{"Female subject", stringPtr("female"), true, "she"},
		{"Female possessive", stringPtr("female"), false, "her"},
		{"Nonbinary subject", stringPtr("nonbinary"), true, "they"},
		{"Nonbinary possessive", stringPtr("nonbinary"), false, "their"},
		{"Nil subject", nil, true, "they"},
		{"Nil possessive", nil, false, "their"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			birthday := database.Birthday{
				Name:   "Test",
				Month:  1,
				Day:    1,
				Gender: tt.gender,
			}

			result := birthday.GetPronoun(tt.subjectForm)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func stringPtr(s string) *string {
	return &s
}
