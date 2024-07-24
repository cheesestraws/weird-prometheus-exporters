package fn

func Contains[A any](as []A, f func(A) bool) bool {
	for _, a := range as {
		if f(a) {
			return true
		}
	}
	return false
}

func Count[A any](as []A, f func(A) bool) int {
	var c int
	for _, a := range as {
		if f(a) {
			c++
		}
	}
	return c
}
