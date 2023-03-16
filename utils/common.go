package utils

func Contains[T string|int|float64|rune](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}