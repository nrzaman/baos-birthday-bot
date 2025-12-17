package birthday

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	"github.com/nrzaman/baos-birthday-bot/internal/database"
	"github.com/nrzaman/baos-birthday-bot/internal/interfaces"
	"github.com/nrzaman/baos-birthday-bot/util"
)

// ServiceDB handles birthday-related operations using a database
type ServiceDB struct {
	timeProvider interfaces.TimeProvider
	db           *database.DB
}

// NewServiceDB creates a new database-backed birthday Service
func NewServiceDB(timeProvider interfaces.TimeProvider, db *database.DB) *ServiceDB {
	return &ServiceDB{
		timeProvider: timeProvider,
		db:           db,
	}
}

// IsBirthdayToday checks if the given month and day match today's date
func (s *ServiceDB) IsBirthdayToday(month int, day int) bool {
	return month == int(s.timeProvider.Month()) && day == s.timeProvider.Day()
}

// GetBirthdayMessage generates a birthday message for anyone with a birthday today
func (s *ServiceDB) GetBirthdayMessage() string {
	now := s.timeProvider.Now()
	birthdays, err := s.db.GetBirthdaysByDate(int(now.Month()), now.Day())
	if err != nil {
		fmt.Printf("Error getting birthdays: %v\n", err)
		return ""
	}

	var buffer bytes.Buffer
	caseyHandled := false
	for _, birthday := range birthdays {
		if birthday.Name == "Casey" && s.IsBirthdayToday(1, 6) && !caseyHandled {
			// Special handling for Casey on January 6th only
			buffer.WriteString("Today is the anniversary of the **Capitol Riots**. Nothing else special happened today.\n")
			caseyHandled = true
		} else {
			// Normal birthday message for everyone else (including Casey on non-1/6 days)
			pronoun := birthday.GetPronoun(false) // possessive form
			buffer.WriteString(fmt.Sprintf("Today is **%s's birthday**! 🎉 Please wish %s a happy birthday! 🎂\n",
				birthday.Name, pronoun))
		}
	}

	return buffer.String()
}

// ListCurrentMonthBirthdays returns a string listing all birthdays in the current month
func (s *ServiceDB) ListCurrentMonthBirthdays() string {
	currentMonth := int(s.timeProvider.Month())

	birthdays, err := s.db.GetBirthdaysByMonth(currentMonth)
	if err != nil {
		fmt.Printf("Error getting birthdays: %v\n", err)
		return ""
	}

	var buffer bytes.Buffer
	for _, birthday := range birthdays {
		month := time.Month(birthday.Month)
		buffer.WriteString(fmt.Sprintf("**%s Birthdays:**\n\n• %s, %s %s\n",
			month.String(),
			birthday.Name,
			month.String(),
			strconv.Itoa(birthday.Day)))
	}

	return buffer.String()
}

// ListAllBirthdays returns a string listing all birthdays
func (s *ServiceDB) ListAllBirthdays() string {
	birthdays, err := s.db.GetAllBirthdays()
	if err != nil {
		fmt.Printf("Error getting birthdays: %v\n", err)
		return ""
	}

	var buffer bytes.Buffer
	buffer.WriteString("**All Birthdays:**\n\n")
	for _, birthday := range birthdays {
		month := time.Month(birthday.Month)
		buffer.WriteString(fmt.Sprintf("• %s, %s %s\n",
			birthday.Name,
			month.String(),
			strconv.Itoa(birthday.Day)))
	}

	return buffer.String()
}

// AddBirthday adds a new birthday
func (s *ServiceDB) AddBirthday(name string, month, day int, gender *string) error {
	return s.db.AddBirthday(name, month, day, gender, nil)
}

// RemoveBirthday removes a birthday
func (s *ServiceDB) RemoveBirthday(name string) error {
	return s.db.DeleteBirthday(name)
}

// GetBirthdaysToday returns all birthdays happening today
func (s *ServiceDB) GetBirthdaysToday() ([]database.Birthday, error) {
	now := s.timeProvider.Now()
	return s.db.GetBirthdaysByDate(int(now.Month()), now.Day())
}

// GetBirthdays returns all birthdays in util.People format for compatibility
func (s *ServiceDB) GetBirthdays() util.People {
	birthdays, err := s.db.GetAllBirthdays()
	if err != nil {
		fmt.Printf("Error getting birthdays: %v\n", err)
		return util.People{People: []util.Person{}}
	}

	people := make([]util.Person, 0, len(birthdays))
	for _, birthday := range birthdays {
		people = append(people, util.Person{
			Name: birthday.Name,
			Birthday: util.Birthday{
				Month: birthday.Month,
				Day:   birthday.Day,
			},
			Gender: birthday.Gender,
		})
	}

	return util.People{People: people}
}
