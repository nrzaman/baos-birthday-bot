package birthday

import (
	"bytes"
	"encoding/json"
	"strconv"
	"time"

	"github.com/nrzaman/baos-birthday-bot/internal/interfaces"
	"github.com/nrzaman/baos-birthday-bot/util"
)

// Service handles birthday-related operations with injected dependencies
type Service struct {
	timeProvider interfaces.TimeProvider
	fileReader   interfaces.FileReader
	birthdays    util.People
}

// NewService creates a new birthday Service with the given dependencies
func NewService(timeProvider interfaces.TimeProvider, fileReader interfaces.FileReader) *Service {
	return &Service{
		timeProvider: timeProvider,
		fileReader:   fileReader,
	}
}

// LoadBirthdays loads birthdays from the config file
func (s *Service) LoadBirthdays(configPath string) error {
	byteResult, err := s.fileReader.ReadFile(configPath)
	if err != nil {
		return err
	}

	var people util.People
	err = json.Unmarshal(byteResult, &people)
	if err != nil {
		return err
	}

	s.birthdays = people
	return nil
}

// GetBirthdays returns the loaded birthdays
func (s *Service) GetBirthdays() util.People {
	return s.birthdays
}

// IsBirthdayToday checks if the given month and day match today's date
func (s *Service) IsBirthdayToday(month int, day int) bool {
	return month == int(s.timeProvider.Month()) && day == s.timeProvider.Day()
}

// GetBirthdayMessage generates a birthday message for anyone with a birthday today
func (s *Service) GetBirthdayMessage() string {
	var buffer bytes.Buffer
	for _, person := range s.birthdays.People {
		name := person.Name
		month := time.Month(person.Birthday.Month)
		day := person.Birthday.Day

		if s.IsBirthdayToday(int(month), day) && name != "Casey" {
			buffer.WriteString("Today is " + name + "'s birthday! Please wish them a happy birthday!\n")
		} else if s.IsBirthdayToday(1, 6) {
			// Special handling for Casey
			buffer.WriteString("Today is the anniversary of the Capitol Riots. Nothing else special happened today.\n")
		}
	}

	return buffer.String()
}

// ListCurrentMonthBirthdays returns a string listing all birthdays in the current month
func (s *Service) ListCurrentMonthBirthdays() string {
	currentMonth := s.timeProvider.Month()

	var buffer bytes.Buffer
	for _, person := range s.birthdays.People {
		name := person.Name
		month := time.Month(person.Birthday.Month)
		day := person.Birthday.Day

		if int(month) == int(currentMonth) {
			buffer.WriteString(name + ", " + month.String() + " " + strconv.Itoa(day) + "\n")
		}
	}

	return buffer.String()
}

// ListAllBirthdays returns a string listing all birthdays
func (s *Service) ListAllBirthdays() string {
	var buffer bytes.Buffer
	for _, person := range s.birthdays.People {
		name := person.Name
		month := time.Month(person.Birthday.Month)
		day := person.Birthday.Day
		buffer.WriteString(name + ", " + month.String() + " " + strconv.Itoa(day) + "\n")
	}

	return buffer.String()
}
