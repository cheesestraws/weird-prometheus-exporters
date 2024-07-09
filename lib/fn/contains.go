package fn

func Contains[A any](as []A, f func(A) bool) bool {
	for _, a := range as {
		if f(a) {
			return true
		}
	}
	return false
}
