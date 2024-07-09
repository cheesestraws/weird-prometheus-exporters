package fn

func Map[A any, B any](as []A, f func(A) B) []B {
	var accum []B
	for _, a := range as {
		accum = append(accum, f(a))
	}
	return accum
}

func Mapkeymap[A comparable, B any, C any](as map[A]B, f func (A) C) []C {
	var accum []C
	for k := range as {
		accum = append(accum, f(k))
	}
	return accum
}