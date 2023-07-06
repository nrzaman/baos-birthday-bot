package birthdayUtil

import "time"

func IsBirthdayCurrentDay(month int, day int) bool {
	now := time.Now()
	currentMonth := now.Month()
	currentDay := now.Day()

	return month == int(currentMonth) && day == currentDay
}
