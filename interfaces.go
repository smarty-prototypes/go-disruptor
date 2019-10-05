package disruptor

type Consumer interface {
	Consume(lower, upper int64)
}

type Barrier interface {
	Load() int64
}

type WaitStrategy interface {
	Gate(int64)
	Idle(int64)
}

type Sequencer interface {
	Reserve(count int64) int64
	Commit(lower, upper int64)
}

type ListenCloser interface {
	Listen()
	Close() error
}
