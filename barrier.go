package main

func (this Barrier) Load() int64 {
	minimum := MaxSequenceValue

	for i := 0; i < len(this); i++ {
		cursor := this[i].Load()
		if cursor < minimum {
			minimum = cursor
		}
	}

	return minimum
}

func NewBarrier(upstream ...*Sequence) Barrier {
	this := Barrier{}
	for i := 0; i < len(upstream); i++ {
		this = append(this, upstream[i])
	}
	return this
}

type Barrier []*Sequence
