package fn

func Filter[A any](as []A, f func(A) bool) []A {
	var accum []A

	for _, a := range as {
		if f(a) {
			accum = append(accum, a)
		}
	}

	return accum
}
