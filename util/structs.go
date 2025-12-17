package util

// Birthday A struct containing the month and day of a birthday.
type Birthday struct {
	Month int `json:"Month"`
	Day   int `json:"Day"`
}

// Person A struct containing the person's first name and birthday.
type Person struct {
	Name     string   `json:"Name"`
	Birthday Birthday `json:"Birthday"`
	Gender   *string  `json:"Gender"`
}

// People A struct containing an array of persons (includes their first name
// and birthday).
type People struct {
	People []Person `json:"Birthdays"`
}
