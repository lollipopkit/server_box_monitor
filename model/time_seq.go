package model

type TimeSequence[T any] struct {
	Old *T
	New *T
}

func (ts *TimeSequence[T]) Update(t *T) {
	ts.Old = ts.New
	ts.New = t
}
