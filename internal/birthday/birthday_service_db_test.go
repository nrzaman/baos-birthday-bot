package birthday_test

import (
	"strings"
	"testing"
	"time"

	"github.com/nrzaman/baos-birthday-bot/internal/birthday"
	"github.com/nrzaman/baos-birthday-bot/internal/database"
)

// MockTimeProvider for testing
type MockTimeProvider struct {
	CurrentTime time.Time
}

func (m *MockTimeProvider) Now() time.Time {
	return m.CurrentTime
}

func (m *MockTimeProvider) Month() time.Month {
	return m.CurrentTime.Month()
}

func (m *MockTimeProvider) Day() int {
	return m.CurrentTime.Day()
}

// Helper function to create test database
func setupTestDB(t *testing.T) *database.DB {
	db, err := database.New(":memory:")
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	return db
}

// Helper function to add test birthdays
func addTestBirthday(t *testing.T, db *database.DB, name string, month, day int, gender *string) {
	err := db.AddBirthday(name, month, day, gender, nil)
	if err != nil {
		t.Fatalf("Failed to add test birthday: %v", err)
	}
}

func TestGetBirthdayMessage_NoBirthdays(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.GetBirthdayMessage()

	// Assert
	if message != "" {
		t.Errorf("Expected empty message when no birthdays, got: %q", message)
	}
}

func TestGetBirthdayMessage_SingleBirthdayMale(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	male := "male"
	addTestBirthday(t, db, "John", 3, 15, &male)

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.GetBirthdayMessage()

	// Assert
	if !strings.Contains(message, "John") {
		t.Errorf("Expected message to contain 'John', got: %q", message)
	}
	if !strings.Contains(message, "his") {
		t.Errorf("Expected message to use 'his' pronoun for male, got: %q", message)
	}
	if !strings.Contains(message, "🎉") {
		t.Errorf("Expected message to contain celebration emoji, got: %q", message)
	}
}

func TestGetBirthdayMessage_SingleBirthdayFemale(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	female := "female"
	addTestBirthday(t, db, "Alice", 3, 15, &female)

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.GetBirthdayMessage()

	// Assert
	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected message to contain 'Alice', got: %q", message)
	}
	if !strings.Contains(message, "her") {
		t.Errorf("Expected message to use 'her' pronoun for female, got: %q", message)
	}
}

func TestGetBirthdayMessage_SingleBirthdayNonbinary(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	nonbinary := "nonbinary"
	addTestBirthday(t, db, "Taylor", 3, 15, &nonbinary)

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.GetBirthdayMessage()

	// Assert
	if !strings.Contains(message, "Taylor") {
		t.Errorf("Expected message to contain 'Taylor', got: %q", message)
	}
	if !strings.Contains(message, "their") {
		t.Errorf("Expected message to use 'their' pronoun for nonbinary, got: %q", message)
	}
}

func TestGetBirthdayMessage_CaseyOnJanuary6th(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	male := "male"
	addTestBirthday(t, db, "Casey", 1, 6, &male)

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 1, 6, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.GetBirthdayMessage()

	// Assert
	if !strings.Contains(message, "Capitol Riots") {
		t.Errorf("Expected Casey's special message on 1/6, got: %q", message)
	}
	if strings.Contains(message, "Casey's birthday") {
		t.Errorf("Expected no normal birthday message for Casey on 1/6, got: %q", message)
	}
	if !strings.Contains(message, "Nothing else special happened today") {
		t.Errorf("Expected full Capitol Riots message, got: %q", message)
	}
}

func TestGetBirthdayMessage_CaseyOnOtherDay(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	male := "male"
	// Casey's birthday is actually 1/6, but we're testing a different day
	addTestBirthday(t, db, "Casey", 3, 15, &male)

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.GetBirthdayMessage()

	// Assert
	if strings.Contains(message, "Capitol Riots") {
		t.Errorf("Expected normal birthday message for Casey on non-1/6 day, got: %q", message)
	}
	if !strings.Contains(message, "Casey") {
		t.Errorf("Expected message to contain Casey's name, got: %q", message)
	}
}

func TestGetBirthdayMessage_January6thWithOtherPeople(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	male := "male"
	female := "female"
	addTestBirthday(t, db, "Casey", 1, 6, &male)
	addTestBirthday(t, db, "Alice", 1, 6, &female)

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 1, 6, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.GetBirthdayMessage()

	// Assert
	if !strings.Contains(message, "Capitol Riots") {
		t.Errorf("Expected Casey's special message on 1/6, got: %q", message)
	}
	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected Alice to still get birthday message on 1/6, got: %q", message)
	}
	// Count occurrences of "Capitol Riots" to ensure it only appears once
	count := strings.Count(message, "Capitol Riots")
	if count != 1 {
		t.Errorf("Expected Capitol Riots message to appear exactly once, appeared %d times", count)
	}
}

func TestGetBirthdayMessage_MultipleBirthdays(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	male := "male"
	female := "female"
	addTestBirthday(t, db, "John", 3, 15, &male)
	addTestBirthday(t, db, "Alice", 3, 15, &female)

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.GetBirthdayMessage()

	// Assert
	if !strings.Contains(message, "John") {
		t.Errorf("Expected message to contain 'John', got: %q", message)
	}
	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected message to contain 'Alice', got: %q", message)
	}
	if !strings.Contains(message, "his") {
		t.Errorf("Expected message to use 'his' for John, got: %q", message)
	}
	if !strings.Contains(message, "her") {
		t.Errorf("Expected message to use 'her' for Alice, got: %q", message)
	}
}

func TestIsBirthdayToday_True(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	result := service.IsBirthdayToday(3, 15)

	// Assert
	if !result {
		t.Error("Expected IsBirthdayToday to return true for matching date")
	}
}

func TestIsBirthdayToday_False(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 15, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	result := service.IsBirthdayToday(3, 16)

	// Assert
	if result {
		t.Error("Expected IsBirthdayToday to return false for non-matching date")
	}
}

func TestListCurrentMonthBirthdays_WithBirthdays(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	male := "male"
	female := "female"
	addTestBirthday(t, db, "John", 3, 15, &male)
	addTestBirthday(t, db, "Alice", 3, 20, &female)
	addTestBirthday(t, db, "Bob", 4, 10, &male) // Different month

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.ListCurrentMonthBirthdays()

	// Assert
	if !strings.Contains(message, "John") {
		t.Errorf("Expected message to contain 'John', got: %q", message)
	}
	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected message to contain 'Alice', got: %q", message)
	}
	if strings.Contains(message, "Bob") {
		t.Errorf("Expected message to NOT contain 'Bob' (different month), got: %q", message)
	}
	if !strings.Contains(message, "March") {
		t.Errorf("Expected message to contain month name, got: %q", message)
	}
}

func TestListCurrentMonthBirthdays_NoBirthdays(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.ListCurrentMonthBirthdays()

	// Assert
	if message != "" {
		t.Errorf("Expected empty message when no birthdays in current month, got: %q", message)
	}
}

func TestListAllBirthdays_MultipleBirthdays(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	male := "male"
	female := "female"
	addTestBirthday(t, db, "John", 3, 15, &male)
	addTestBirthday(t, db, "Alice", 6, 20, &female)
	addTestBirthday(t, db, "Bob", 12, 25, &male)

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.ListAllBirthdays()

	// Assert
	if !strings.Contains(message, "All Birthdays") {
		t.Errorf("Expected header 'All Birthdays', got: %q", message)
	}
	if !strings.Contains(message, "John") {
		t.Errorf("Expected message to contain 'John', got: %q", message)
	}
	if !strings.Contains(message, "Alice") {
		t.Errorf("Expected message to contain 'Alice', got: %q", message)
	}
	if !strings.Contains(message, "Bob") {
		t.Errorf("Expected message to contain 'Bob', got: %q", message)
	}
	if !strings.Contains(message, "March") {
		t.Errorf("Expected message to contain 'March', got: %q", message)
	}
	if !strings.Contains(message, "June") {
		t.Errorf("Expected message to contain 'June', got: %q", message)
	}
	if !strings.Contains(message, "December") {
		t.Errorf("Expected message to contain 'December', got: %q", message)
	}
}

func TestListAllBirthdays_NoBirthdays(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	message := service.ListAllBirthdays()

	// Assert
	if !strings.Contains(message, "All Birthdays") {
		t.Errorf("Expected header even with no birthdays, got: %q", message)
	}
}

func TestGetBirthdays_ReturnsCorrectFormat(t *testing.T) {
	// Arrange
	db := setupTestDB(t)
	defer db.Close()

	male := "male"
	addTestBirthday(t, db, "John", 3, 15, &male)

	timeProvider := &MockTimeProvider{
		CurrentTime: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
	}
	service := birthday.NewServiceDB(timeProvider, db)

	// Act
	people := service.GetBirthdays()

	// Assert
	if len(people.People) != 1 {
		t.Fatalf("Expected 1 person, got %d", len(people.People))
	}
	if people.People[0].Name != "John" {
		t.Errorf("Expected name 'John', got %q", people.People[0].Name)
	}
	if people.People[0].Birthday.Month != 3 {
		t.Errorf("Expected month 3, got %d", people.People[0].Birthday.Month)
	}
	if people.People[0].Birthday.Day != 15 {
		t.Errorf("Expected day 15, got %d", people.People[0].Birthday.Day)
	}
	if people.People[0].Gender == nil || *people.People[0].Gender != "male" {
		t.Errorf("Expected gender 'male', got %v", people.People[0].Gender)
	}
}
