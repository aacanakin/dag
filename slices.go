package dag

// filter filters out a slice given test func
func filter[T any](slice []T, test func(T) bool) (ret []T) {
	for _, item := range slice {
		if test(item) {
			ret = append(ret, item)
		}
	}
	return
}

func index[T any](slice []T, test func(T) bool) (index int) {
	for index, item := range slice {
		if test(item) {
			return index
		}
	}
	return -1
}

// includes checks if a slice includes a search item
func includes[T comparable](slice []T, search T) bool {
	idx := index(slice, func(v T) bool {
		return search == v
	})
	return idx >= 0
}

// some returns true when a test func satisfies a slice once
func some[T any](slice []T, test func(T) bool) bool {
	for _, item := range slice {
		if test(item) {
			return true
		}
	}
	return false
}
