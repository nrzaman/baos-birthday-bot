package birthday

import "github.com/nrzaman/baos-birthday-bot/util"

// BirthdayService defines the interface for birthday-related operations
type BirthdayService interface {
	// IsBirthdayToday checks if the given month and day match today's date
	IsBirthdayToday(month int, day int) bool

	// GetBirthdayMessage generates a birthday message for anyone with a birthday today
	GetBirthdayMessage() string

	// ListCurrentMonthBirthdays returns a string listing all birthdays in the current month
	ListCurrentMonthBirthdays() string

	// ListAllBirthdays returns a string listing all birthdays
	ListAllBirthdays() string

	// GetBirthdays returns all birthdays (for compatibility with existing code)
	GetBirthdays() util.People
}
