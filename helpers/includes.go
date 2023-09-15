package helpers

// iterates through given array and returns a bool of if
// val was found in the array
func Includes[T string | int](array []T, val T) bool {
	for _, item := range array {
		if val == item {
			return true
		}
	}
	return false
}
