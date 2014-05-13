package main

func (this Barrier) Load() int64 {
	minimum := MaxSequenceValue
	length := this.length
	upstream := this.upstream

	for i := 0; i < length; i++ {
		cursor := upstream[i].Load()
		if cursor < minimum {
			minimum = cursor
		}
	}

	return minimum
}

func NewBarrier(upstream ...*Sequence) Barrier {
	length := len(upstream)
	target := make([]*Sequence, length, length)
	copy(target, upstream)
	return Barrier{length: length, upstream: target}
}

type Barrier struct {
	length   int
	upstream []*Sequence
}
