package benchmarks

import "runtime"

const (
	RingBufferSize   = 1024 * 64
	RingBufferMask   = RingBufferSize - 1
	ReserveOne       = 1
	ReserveMany      = 16
	ReserveManyDelta = ReserveMany - 1
)

func init() {
	runtime.GOMAXPROCS(2)
}
