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