package helpers

import (
	"strconv"
)

// check if the given id appears to be a number
func IsNumber(id string) bool {
	if _, err := strconv.Atoi(id); err == nil {
		return true
	}
	return false
}
