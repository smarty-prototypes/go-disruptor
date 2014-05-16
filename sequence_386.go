package disruptor

import "sync/atomic"

func (this *Sequence) Store(value int64) {
	atomic.StoreInt64(&(*this)[SequencePayloadIndex], value)
}
func (this *Sequence) Load() int64 {
	return atomic.LoadInt64(&(*this)[SequencePayloadIndex])
}
