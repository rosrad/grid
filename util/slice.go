package util

func Contains(s []int, d int) bool {
	for _, i := range s {
		if d == i {
			return true
		}
	}
	return false
}
