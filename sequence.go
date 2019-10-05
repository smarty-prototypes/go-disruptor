package disruptor

import "sync/atomic"

type Sequence [8]int64 // prevent false sharing of the sequence cursor by padding the CPU cache line with 64 *bytes* of data.

func NewCursor() *Sequence {
	var this Sequence
	this[0] = defaultSequenceValue
	return &this
}

func (this *Sequence) Store(sequence int64) { atomic.StoreInt64(&this[0], sequence) }
func (this *Sequence) Load() int64          { return atomic.LoadInt64(&this[0]) }

const defaultSequenceValue = -1
