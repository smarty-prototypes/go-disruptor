package disruptor

import "sync/atomic"

type Cursor [cpuCacheLinePadding]int64

func NewCursor() *Cursor {
	this := &Cursor{}
	this.Store(InitialCursorSequenceValue)
	return this
}

func (this *Cursor) Store(sequence int64) { atomic.StoreInt64(&this[0], sequence) }
func (this *Cursor) Load() int64          { return atomic.LoadInt64(&this[0]) }

const cpuCacheLinePadding = 8
