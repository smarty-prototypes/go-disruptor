package disruptor

type Disruptor interface {
	Writer
	ListenCloser
}

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

const ErrReservationSize = -1
