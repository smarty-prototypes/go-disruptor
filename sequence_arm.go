package main

import "sync/atomic"

// TODO: looks like there's one more race condition in the ARM processor.

func (this *Sequence) Store(value int64) {
	atomic.StoreInt64(&(*this)[SequencePayloadIndex], value)
}
func (this *Sequence) Load() int64 {
	return atomic.LoadInt64(&(*this)[SequencePayloadIndex])
}

const FillCPUCacheLine uint8 = 8 // TODO: how big is the cache line for ARM?
