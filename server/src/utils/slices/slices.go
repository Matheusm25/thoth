package SliceUtils

func Contains[T comparable](slice []T, target T) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}
	return false
}

func Filter[T any](slice []T, filter func(T) bool) []T {
	var filtered []T
	for _, v := range slice {
		if filter(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}
