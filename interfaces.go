package disruptor

import "errors"

type Handler interface {
	Handle(lower, upper int64)
}

type sequenceBarrier interface {
	Load() int64
}

type WaitStrategy interface {
	Gate(int64)
	Idle(int64)
}

type Writer interface {
	Reserve(count int64) int64
	Commit(lower, upper int64)
}

type ListenCloser interface {
	Listen()
	Close() error
}

var ErrMinimumReservationSize = errors.New("the minimum reservation size is 1 slot")
