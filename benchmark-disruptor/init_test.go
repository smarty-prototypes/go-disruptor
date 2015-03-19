package benchmarks

import "time"

const (
	RingBufferSize   = 64
	RingBufferMask   = RingBufferSize - 1
	ReserveOne       = 1
	ReserveMany      = 16
	ReserveManyDelta = ReserveMany - 1
	DisruptorCleanup = time.Millisecond * 10
)

var ringBuffer = [RingBufferSize]int64{}
