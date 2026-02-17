package disruptor

import (
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

type HandleWaitStrategy interface {
	Gate(int64)
	Idle(int64)
}
type ReserveWaitStrategy interface {
	Wait()
}

type Handler interface {
	Handle(lowerSequence, upperSequence int64)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// The Sequencer tracks the state of a given "writer" or producer to the ring buffer. It is the heart of the Disruptor
// pattern. When a caller desires to push events to the ring buffer, the caller or product must first Reserve the
// desired slots on the ring buffer using a Sequencer. After obtaining a reservation on certain slots of the ring
// buffer, the caller must write the any data to the reserved slots and then indicate the completion of the
// operation(s) by calling Commit which makes those slots available to any downstream Handler instances which then
// handle or otherwise consume events from the ring buffer on different goroutines.
type Sequencer interface {

	// Reserve claims the desired number of slots in the ring buffer for the caller. When those slots become available
	// because any configured Handlers have properly processed all necessary data in those slots, the Sequencer returns
	// the uppermost or highest sequence of slot claimed and reserved for the caller.
	//
	// The lower-bound sequence in the ring buffer is obtained by subtracting the specified number of slots from the
	// uppermost sequence returned. If the number of desired slots is larger than the capacity of the ring buffer,
	// ErrReservationSize is returned.
	//
	// Each successful call to Reserve should *always* be followed by a single call to Commit.
	Reserve(slots int64) (upperSequence int64)

	// Commit indicates to the sequencer that the previously claimed slots in the ring buffer have been written to
	// successfully and to make the data available to any configured Handler instances to process.
	//
	// Each successful call to Commit should *always* be preceded by a single successful call to Reserve.
	Commit(lowerSequence, upperSequence int64)
}

const (
	// ErrReservationSize indicates that the reservation requested is nonsensical, e.g. lower > upper OR that the
	// desired reservation size exceeds the capacity of the ring buffer altogether.
	ErrReservationSize = -1
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type atomicSequence = *atomic.Int64

type sequenceBarrier interface {
	Load(int64) int64
}

type atomicBarrier struct{ sequence atomicSequence }

func newAtomicBarrier(sequence atomicSequence) *atomicBarrier {
	return &atomicBarrier{sequence: sequence}
}
func (this *atomicBarrier) Load(_ int64) int64 { return this.sequence.Load() }

const defaultSequenceValue = -1
