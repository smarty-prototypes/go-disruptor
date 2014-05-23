package disruptor

import "sync/atomic"

func (this *Cursor) Store(sequence int64) {
	atomic.StoreInt64(&this.sequence, sequence)
}
func (this *Cursor) Load() int64 {
	return atomic.LoadInt64(&this.sequence)
}
func (this *Cursor) LoadBarrier(next int64) int64 {
	return atomic.LoadInt64(&this.sequence)
}
