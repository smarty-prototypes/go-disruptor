package main

import "sync/atomic"

type Sequence [FillCPUCacheLine]int64

func (this *Sequence) Store(value int64) {
	// TODO: this needs build tags for i386, arm, and amd64 because of torn writes
	// atomic.StoreInt64(&(*this)[SequencePayloadIndex], value)
	(*this)[SequencePayloadIndex] = value
}
func (this *Sequence) Load() int64 {
	// TODO: this needs build tags for i386, arm, and amd64 because of torn writes

	// interestingly, running atomic.Load (but normal/regular store) on x86_64
	// makes things FASTER, e.g. 700ms per 100 million operations instead of 730ms.
	// One theory is that the golang scheduler ppears to try to make things efficient
	// by running them on a single core, so atomic makes the routine run slower thus
	// the scheduler keeps things running on multiple cores.

	return atomic.LoadInt64(&(*this)[SequencePayloadIndex])
	// return (*this)[SequencePayloadIndex]
}

func NewSequence() *Sequence {
	return &Sequence{InitialSequenceValue}
}

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	SequencePayloadIndex uint8 = 0
	FillCPUCacheLine     uint8 = 8 // FUTURE: use build tags for i386, amd64, and ARM-v4,5,6,7,8 processors
)
