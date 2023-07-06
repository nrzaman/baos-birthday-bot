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
}

// People A struct containing an array of persons (includes their first name
// and birthday).
type People struct {
	People []Person `json:"Birthdays"`
}

// Channel A struct containing the channel name and ID.
type Channel struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
}

// Channels A struct containing an array of channels.
type Channels struct {
	Channel []Channel `json:"Channels"`
}
