package sliceutils

func SliceContains(slice []string, needle string) bool {
	for _, elem := range slice {
		if elem == needle {
			return true
		}
	}
	return false
}
