package birthdayUtil

import "log"

func Check(e error) {
	if e != nil {
		log.Fatal("Error: ", e)
	}
}
