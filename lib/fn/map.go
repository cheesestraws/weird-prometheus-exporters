package fn

func Map[A any, B any](as []A, f func(A) B) []B {
	var accum []B
	for _, a := range as {
		accum = append(accum, f(a))
	}
	return accum
}

func Mapmap[A comparable, B any, C comparable, D any](bs map[A]B,
	f func(A, B) (C, D)) map[C]D {
	
	m := make(map[C]D)
	for a, b := range bs {
		c, d := f(a, b)
		m[c] = d
	}
	
	return m
}

func SimpleMapmap[A comparable, B any, C any](bs map[A]B,
	f func(B) C) map[A]C {
	
	m := make(map[A]C)
	for a, b := range bs {
		m[a] = f(b)
	}	
	
	return m
}


func Mapkeymap[A comparable, B any, C any](as map[A]B, f func(A) C) []C {
	var accum []C
	for k := range as {
		accum = append(accum, f(k))
	}
	return accum
}

func Map2mapkey[A any, B comparable, C any](as []A, f func(A) (B, C)) map[B]C {
	accum := make(map[B]C)

	for _, a := range as {
		b, c := f(a)
		accum[b] = c
	}

	return accum
}

func Dedupe[A comparable](as []A) []A {
	m := Map2mapkey(as, func(a A) (A, struct{}) {
		return a, struct{}{}
	})

	return Mapkeymap(m, Id[A])
}
