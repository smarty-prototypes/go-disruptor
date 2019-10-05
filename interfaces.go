package disruptor

import "io"

type Listener interface {
	Listen()
}
type ListenCloser interface {
	Listener
	io.Closer
}

type Consumer interface {
	Consume(lower, upper int64)
}

type Barrier interface {
	Load() int64
}

type WaitStrategy interface {
	Gate(int)
	Idle(int)
}

type Sequencer interface {
	Reserve(count int64) int64
	Commit(lower, upper int64)
}

const (
	MaxCursorSequenceValue     int64 = (1 << 63) - 1
	InitialCursorSequenceValue int64 = -1
)
