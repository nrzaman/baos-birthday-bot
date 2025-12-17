package birthday_test

import (
	"testing"
	"time"

	"github.com/nrzaman/baos-birthday-bot/internal/birthday"
)

// MockTimeProvider is a mock implementation of TimeProvider for testing
type MockTimeProvider struct {
	CurrentTime time.Time
}

// Get current time
func (m *MockTimeProvider) Now() time.Time {
	return m.CurrentTime
}

func (m *MockTimeProvider) Month() time.Month {
	return m.CurrentTime.Month()
}

func (m *MockTimeProvider) Day() int {
	return m.CurrentTime.Day()
}

// MockFileReader is a mock implementation of FileReader for testing
type MockFileReader struct {
	Data map[string][]byte
	Err  error
}

func (m *MockFileReader) ReadFile(path string) ([]byte, error) {
	if m.Err != nil {
		return nil, m.Err
	}
	return m.Data[path], nil
}

func TestIsBirthdayToday(t *testing.T) {
	// Arrange
	mockTime := &MockTimeProvider{
		CurrentTime: time.Date(2024, time.March, 15, 10, 0, 0, 0, time.UTC),
	}
	mockFileReader := &MockFileReader{Data: map[string][]byte{}}

	service := birthday.NewService(mockTime, mockFileReader)

	tests := []struct {
		name     string
		month    int
		day      int
		expected bool
	}{
		{"Same month and day", 3, 15, true},
		{"Different month", 4, 15, false},
		{"Different day", 3, 16, false},
		{"Completely different", 12, 25, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := service.IsBirthdayToday(tt.month, tt.day)

			// Assert
			if result != tt.expected {
				t.Errorf("IsBirthdayToday(%d, %d) = %v; want %v", tt.month, tt.day, result, tt.expected)
			}
		})
	}
}

func TestGetBirthdayMessage(t *testing.T) {
	// Arrange - Set current date to March 15
	mockTime := &MockTimeProvider{
		CurrentTime: time.Date(2024, time.March, 15, 10, 0, 0, 0, time.UTC),
	}

	birthdaysJSON := `{
		"Birthdays": [
			{"Name": "Alice", "Birthday": {"Month": 3, "Day": 15}},
			{"Name": "Bob", "Birthday": {"Month": 3, "Day": 16}},
			{"Name": "Charlie", "Birthday": {"Month": 4, "Day": 15}}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/birthdays.json": []byte(birthdaysJSON),
		},
	}

	service := birthday.NewService(mockTime, mockFileReader)
	err := service.LoadBirthdays("./config/birthdays.json")
	if err != nil {
		t.Fatalf("Failed to load birthdays: %v", err)
	}

	// Act
	message := service.GetBirthdayMessage()

	// Assert
	expected := "Today is Alice's birthday! Please wish them a happy birthday!\n"
	if message != expected {
		t.Errorf("GetBirthdayMessage() = %q; want %q", message, expected)
	}
}

func TestListCurrentMonthBirthdays(t *testing.T) {
	// Arrange - Set current date to March 15
	mockTime := &MockTimeProvider{
		CurrentTime: time.Date(2024, time.March, 15, 10, 0, 0, 0, time.UTC),
	}

	birthdaysJSON := `{
		"Birthdays": [
			{"Name": "Alice", "Birthday": {"Month": 3, "Day": 15}},
			{"Name": "Bob", "Birthday": {"Month": 3, "Day": 20}},
			{"Name": "Charlie", "Birthday": {"Month": 4, "Day": 10}}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/birthdays.json": []byte(birthdaysJSON),
		},
	}

	service := birthday.NewService(mockTime, mockFileReader)
	err := service.LoadBirthdays("./config/birthdays.json")
	if err != nil {
		t.Fatalf("Failed to load birthdays: %v", err)
	}

	// Act
	result := service.ListCurrentMonthBirthdays()

	// Assert
	// Should only include March birthdays
	if !contains(result, "Alice") || !contains(result, "Bob") {
		t.Errorf("ListCurrentMonthBirthdays() should include Alice and Bob")
	}
	if contains(result, "Charlie") {
		t.Errorf("ListCurrentMonthBirthdays() should not include Charlie (April birthday)")
	}
}

func TestListAllBirthdays(t *testing.T) {
	// Arrange
	mockTime := &MockTimeProvider{
		CurrentTime: time.Date(2024, time.March, 15, 10, 0, 0, 0, time.UTC),
	}

	birthdaysJSON := `{
		"Birthdays": [
			{"Name": "Alice", "Birthday": {"Month": 3, "Day": 15}},
			{"Name": "Bob", "Birthday": {"Month": 6, "Day": 20}},
			{"Name": "Charlie", "Birthday": {"Month": 12, "Day": 25}}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/birthdays.json": []byte(birthdaysJSON),
		},
	}

	service := birthday.NewService(mockTime, mockFileReader)
	err := service.LoadBirthdays("./config/birthdays.json")
	if err != nil {
		t.Fatalf("Failed to load birthdays: %v", err)
	}

	// Act
	result := service.ListAllBirthdays()

	// Assert
	// Should include all birthdays
	if !contains(result, "Alice") || !contains(result, "Bob") || !contains(result, "Charlie") {
		t.Errorf("ListAllBirthdays() should include all people: Alice, Bob, and Charlie")
	}
}

func TestLoadBirthdays(t *testing.T) {
	// Arrange
	mockTime := &MockTimeProvider{}
	birthdaysJSON := `{
		"Birthdays": [
			{"Name": "Alice", "Birthday": {"Month": 3, "Day": 15}}
		]
	}`

	mockFileReader := &MockFileReader{
		Data: map[string][]byte{
			"./config/birthdays.json": []byte(birthdaysJSON),
		},
	}

	service := birthday.NewService(mockTime, mockFileReader)

	// Act
	err := service.LoadBirthdays("./config/birthdays.json")

	// Assert
	if err != nil {
		t.Errorf("LoadBirthdays() returned error: %v", err)
	}

	birthdays := service.GetBirthdays()
	if len(birthdays.People) != 1 {
		t.Errorf("Expected 1 birthday, got %d", len(birthdays.People))
	}

	if birthdays.People[0].Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", birthdays.People[0].Name)
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
