package birthdayUtil

type Birthday struct {
	Month int `json:"Month"`
	Day   int `json:"Day"`
}

type Person struct {
	Name     string   `json:"Name"`
	Birthday Birthday `json:"Birthday"`
}

type People struct {
	People []Person `json:"Birthdays"`
}
