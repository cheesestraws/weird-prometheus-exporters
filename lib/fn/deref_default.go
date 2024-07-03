package fn

func DerefOrDefault[A any](ptr *A, defaultValue A) A {
	if ptr == nil {
		return defaultValue
	}
	return *ptr
}
