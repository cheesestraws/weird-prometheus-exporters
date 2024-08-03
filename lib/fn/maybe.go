package fn

type Maybe[T any] struct {
	ok bool
	t T
}

func Present[T any](t T) Maybe[T] {
	return Maybe[T]{
		ok: true,
		t: t,
	}
}

func Absent[T any]() Maybe[T] {
	return Maybe[T]{
		ok: false,
	}
}

func (m Maybe[T]) Get() (T, bool) {
	return m.t, m.ok
}

func (m Maybe[T]) Range(f func(t T)) Maybe[T] {
	if m.ok {
		f(m.t)
	}
	return m
}

func (m Maybe[T]) OrElse(f func()) Maybe[T] {
	if !m.ok {
		f()
	}
	return m
}

func (m Maybe[T]) Or(v T) T {
	if m.ok {
		return m.t
	} else {
		return v
	}
}