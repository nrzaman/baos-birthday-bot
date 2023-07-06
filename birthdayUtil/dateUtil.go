package birthdayUtil

import "time"

// IsBirthdayCurrentDay This function checks whether the current month and current day
// match the inputted birthday (returns true if so, false if not)
func IsBirthdayCurrentDay(month int, day int) bool {
	return month == int(GetCurrentMonth()) && day == GetCurrentDay()
}

// GetCurrentMonth Returns the current month.
func GetCurrentMonth() time.Month {
	return time.Now().Month()
}

// GetCurrentDay Returns the current day.
func GetCurrentDay() int {
	return time.Now().Day()
}
