package folder

import "os"

// DoesNotExist returns true if the folder does not exist,
// false if the folder DOES exist
func DoesNotExist(folder string) bool {
	retval := false

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		retval = true
	}

	return retval
}

// DoesExist returns true if the folder exists
// false if the folder DOES NOT exist
func DoesExist(folder string) bool {
	retval := true

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		retval = false
	}

	return retval
}
