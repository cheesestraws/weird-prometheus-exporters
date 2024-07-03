package fn

func Map[A any, B any](as []A, f func(A) B) []B {
	var accum []B
	for _, a := range as {
		accum = append(accum, f(a))
	}
	return accum
}
