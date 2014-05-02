package disruptor

import "sync/atomic"

type Sequence []uint64

func (this Sequence) AtomicLoad() uint64 {
	return atomic.LoadUint64(&this[0])
}
func (this Sequence) Store(value uint64) {
	this[0] = value
}

func NewSequence() Sequence {
	return Sequence(make([]uint64, FillCPUCacheLine))
}

const FillCPUCacheLine = 8
