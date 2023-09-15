package helpers

import "errors"

func RemoveElement(array []string, element string) ([]string, error) {
	newArray := make([]string, 0)
	index := getIndex(array, element)
	if index == -1 {
		return newArray, errors.New("element not in array")
	}
	newArray = append(newArray, array[:index]...)
	return append(newArray, array[index+1:]...), nil

}

func getIndex(array []string, element string) int {
	for i, v := range array {
		if v == element {
			return i
		}
	}
	return -1
}
