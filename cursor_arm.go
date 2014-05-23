package disruptor

import "sync/atomic"

func (this *Cursor) Store(value int64) {
	atomic.StoreInt64(&this.value, value)
}
func (this *Cursor) Load() int64 {
	return atomic.LoadInt64(&this.value)
}
func (this *Cursor) LoadBarrier(current int64) int64 {
	return atomic.LoadInt64(&this.value)
}
