package model

type TimeSequence[T CpuOneTimeStatus | NetworkOneTimeStatus] struct {
	Old *T
	New *T
}

func (ts *TimeSequence[T]) Update(t *T) {
	ts.Old = ts.New
	ts.New = t
}
