package disruptor

import "sync/atomic"

type Cursor struct {
	sequence int64
	padding  [cpuCacheLinePadding]int64
}

func NewCursor() *Cursor {
	return &Cursor{sequence: InitialSequenceValue}
}

func (this *Cursor) Store(sequence int64) { atomic.StoreInt64(&this.sequence, sequence) }
func (this *Cursor) Load() int64          { return atomic.LoadInt64(&this.sequence) }
func (this *Cursor) Read(_ int64) int64   { return atomic.LoadInt64(&this.sequence) }
func (this *Cursor) Swap(previous, updated int64) bool {
	return atomic.CompareAndSwapInt64(&this.sequence, previous, updated)
}

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	cpuCacheLinePadding        = 7
)
