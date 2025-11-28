package disruptor

import "io"

type Disruptor[T any] interface {
	Writers() []Writer
	ListenCloser
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type ListenCloser interface {
	Listener
	io.Closer
}
type Listener interface {
	Listen()
}

type WaitStrategy interface {
	Gate(int64)
	Idle(int64)
}

type Handler interface {
	Handle(lowerSequence, upperSequence int64)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Writer interface {
	Reserve(slots int64) (lowerSequence int64)
	Commit(lowerSequence, upperSequence int64)
}

const ErrReservationSize = -1

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type sequenceBarrier interface {
	Load() int64
}

const defaultCursorValue = -1
