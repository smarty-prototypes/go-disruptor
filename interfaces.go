package disruptor

import (
	"context"
	"io"
	"sync/atomic"
)

type Disruptor interface {
	Sequencers() []Sequencer
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

type Sequencer interface {
	Reserve(ctx context.Context, slots int64) (upperSequence int64)
	Commit(lowerSequence, upperSequence int64)
}

const (
	ErrReservationSize = -1
	ErrContextCanceled = -2
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type atomicSequence = *atomic.Int64

type sequenceBarrier interface {
	Load() int64
}

const defaultSequenceValue = -1
