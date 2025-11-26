package disruptor

import (
	"math"
	"sync/atomic"
)

type multiBarrier []*atomic.Int64

func NewMultiBarrier(sequences ...*atomic.Int64) Barrier {
	if len(sequences) == 1 {
		return sequences[0]
	} else {
		return multiBarrier(sequences)
	}
}

func (this multiBarrier) Load() int64 {
	var minimum int64 = math.MaxInt64

	for _, item := range this {
		if sequence := item.Load(); sequence < minimum {
			minimum = sequence
		}
	}

	return minimum
}
