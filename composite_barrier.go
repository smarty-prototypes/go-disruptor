package disruptor

import "math"

type compositeBarrier []*Sequence

func (this compositeBarrier) Load() int64 {
	var minimum int64 = math.MaxInt64

	for _, item := range this {
		if sequence := item.Load(); sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}
