package fn

func Errmap[A any, B any](as []A, f func(A) (B, error)) ([]B, error) {
	var accum []B

	for _, a := range as {
		b, err := f(a)
		if err != nil {
			return nil, err
		}
		accum = append(accum, b)
	}

	return accum, nil
}

func Errmapmap[A any, B comparable, C any](as []A, f func(A) (B, C, error)) (map[B]C, error) {
	accum := make(map[B]C)

	for _, a := range as {
		b, c, err := f(a)
		if err != nil {
			return nil, err
		}
		accum[b] = c
	}

	return accum, nil
}
