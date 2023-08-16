package utils

func Includes(arr []string, searchStr string) bool {
	found := false
	for _, item := range arr {
		if item == searchStr {
			found = true
			break
		}
	}
	return found

}
