package util

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"time"
)

var Birthdays People

// ExtractBirthdays This function extracts birthdays from the JSON config file and stores them
// to be referenced later.
func ExtractBirthdays() {
	// Open the JSON config file
	content, err := os.Open("./config/birthdays.json")
	Check(err)

	defer func(content *os.File) {
		err := content.Close()
		Check(err)
	}(content)

	// Read all contents
	byteResult, _ := io.ReadAll(content)

	// Create result variable
	var people People

	// Store contents
	err = json.Unmarshal(byteResult, &people)
	Check(err)
	Birthdays = people
}

// BirthdayFinderMessage This function will be called in order to determine whether the current
// day is a birthday and subsequently constructs and returns a birthday message string to be posted on
// the Discord server.
func BirthdayFinderMessage() string {
	var buffer bytes.Buffer
	for _, person := range Birthdays.People {
		// Get each person's name and date of birth
		name := person.Name
		month := time.Month(person.Birthday.Month)
		day := person.Birthday.Day

		// Check the person's birthday against current day
		if IsBirthdayCurrentDay(int(month), day) && name != "Casey" {
			buffer.WriteString("Today is " + name + "'s birthday! Please wish them a happy birthday!\n")
		} else if IsBirthdayCurrentDay(1, 6) {
			// Special handling for Casey
			buffer.WriteString("Today is the anniversary of the Capitol Riots. Nothing else special happened today.\n")
		}
	}

	return buffer.String()
}
