package disruptor

import "math"

type CompositeBarrier []*Sequence

func NewCompositeBarrier(sequences []*Sequence) CompositeBarrier { return sequences }

func (this CompositeBarrier) Load() int64 {
	var minimum int64 = math.MaxInt64

	for _, item := range this {
		if sequence := item.Load(); sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}
