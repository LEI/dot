package sliceutils

// Contains check if a slice contains a given string
func Contains(in []string, s string) bool {
	for _, a := range in {
		if a == s {
			return true
		}
	}
	return false
}

// Intersects check if a slice contains any element from the other one
func Intersects(in []string, list []string) bool {
	for _, a := range in {
		if Contains(list, a) {
			return true
		}
	}
	return false
}
