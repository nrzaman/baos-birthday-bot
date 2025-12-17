package util

import "log"

// Check This function checks for errors and logs if there is an error.
func Check(e error) {
	if e != nil {
		log.Fatal("Error: ", e)
	}
}
