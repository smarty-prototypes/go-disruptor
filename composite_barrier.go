package disruptor

import (
	"math"
	"sync/atomic"
)

type compositeBarrier []*atomic.Int64

func newCompositeBarrier(sequences ...*atomic.Int64) sequenceBarrier {

	if len(sequences) == 0 {
		return compositeBarrier(nil)
	} else if len(sequences) == 1 {
		return sequences[0]
	} else {
		return compositeBarrier(sequences)
	}
}

func (this compositeBarrier) Load() int64 {
	var minimum int64 = math.MaxInt64

	for _, item := range this {
		if sequence := item.Load(); sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}
