package disruptor

import (
	"context"
	"io"
	"sync/atomic"
)

// Disruptor is the top-level container combining a Sequencer (producer API) with a ListenCloser (consumer
// lifecycle). Created via invocations of New(...option).
type Disruptor interface {
	Sequencer
	ListenCloser
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// ListenCloser combines Listener with io.Closer, allowing consumers to be started and stopped.
type ListenCloser interface {
	Listener
	io.Closer
}

// Listener processes events from the ring buffer. Listen blocks the calling goroutine until the listener is closed.
type Listener interface {
	Listen()
}

// WaitStrategy provides pluggable backpressure for both producers and consumers. The default implementation uses
// time.Sleep(time.Nanosecond) for Gate, time.Sleep(time.Millisecond) for Idle, and runtime.Gosched for Reserve.
type WaitStrategy interface {
	// Gate is invoked when data has been committed to the ring buffer by a producer but the prior Handler group has
	// not yet finished processing it. This means new work is imminent and the Listener should wait briefly.
	Gate(int64)
	// Idle is invoked when all meaningful work has been completed and there are no slots with available messages.
	Idle(int64)
	// Reserve is invoked by the Sequencer when there are no available slots in the underlying ring buffer.
	Reserve(int64)
	// TryReserve is invoked by the Sequencer when there are no available slots in the underlying ring buffer.
	TryReserve(context.Context) error
}

// Handler is the consumer callback invoked by a Listener with each batch of available sequences from the ring buffer.
type Handler interface {
	Handle(lowerSequence, upperSequence int64)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// The Sequencer tracks the state of a given "writer" or producer to the ring buffer. It is the heart of the Disruptor
// pattern. When a caller desires to push events to the ring buffer, the caller or producer must first Reserve the
// desired slots on the ring buffer using the Sequencer. After obtaining a reservation on certain slots of the ring
// buffer, the caller must write the desired data to the reserved slots of the associated ring buffer and then
// indicate the completion of the operation(s) by calling Commit which makes the data in the written slots available
// to any downstream Handler instances which then handle or otherwise consume events from the ring buffer on different
// goroutines.
type Sequencer interface {

	// Reserve claims the desired number of slots in the ring buffer for the caller. When those slots become available
	// because any configured Handlers have properly processed all necessary data in those slots, the Sequencer returns
	// the uppermost or highest sequence of the slots claimed and reserved for the caller.
	//
	// The lower-bound sequence in the ring buffer is obtained by subtracting the specified number of slots from the
	// uppermost sequence returned. If the number of desired slots is larger than the capacity of the ring buffer,
	// ErrReservationSize is returned.
	//
	// Each successful call to Reserve should *always* be followed by a single call to Commit.
	Reserve(slots uint32) (upperSequence int64)

	// TryReserve behaves like Reserve but can be canceled via the provided context.Context. If the context is
	// canceled before the slots can be successfully claimed, ErrContextCanceled is returned. For the shared
	// Sequencer, which is capable of being shared by multiple goroutines, this uses a CAS-based loop instead of
	// atomic Add, which is slower under contention but allows cancellation.
	//
	// Each successful call to TryReserve should *always* be followed by a single call to Commit.
	TryReserve(ctx context.Context, slots uint32) (upperSequence int64)

	// Commit indicates to the Sequencer that the previously claimed slots in the ring buffer have been written to
	// successfully and that the data is now available to any configured Handler instances to process.
	//
	// Each successful call to Commit should *always* be preceded by a single successful call to Reserve.
	Commit(lowerSequence, upperSequence int64)
}

const (
	// ErrReservationSize indicates that the reservation requested is nonsensical, e.g. lower > upper OR that the
	// desired reservation size exceeds the capacity of the ring buffer altogether.
	ErrReservationSize = -1
	// ErrContextCanceled indicates that the reservation request failed because the provided context.Context has
	// been canceled or otherwise timed out.
	ErrContextCanceled = -2

	// spinMask controls how often the context is checked in TryReserve slow paths.
	// Must be 2^n-1. Context is checked every spinMask+1 iterations.
	spinMask = 1024*16 - 1
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// atomicSequence is a cache-line-padded atomic int64 used to track sequence positions without false sharing. The 56
// bytes of padding on each side ensure the embedded atomic.Int64 occupies its own cache line.
type atomicSequence struct {
	_ [7]int64 // 56B left padding
	atomic.Int64
	_ [7]int64 // 56B right padding
}

// sequenceBarrier abstracts reading the committed or handled position from one or more sequences. Load returns the
// highest sequence that is safe to read up to from the given lower bound.
type sequenceBarrier interface {
	Load(int64) int64
}

// atomicBarrier is a sequenceBarrier backed by a single atomicSequence. Used when there is exactly one upstream
// sequence to track, avoiding the iteration overhead of compositeBarrier.
type atomicBarrier struct{ sequence *atomicSequence }

func newAtomicBarrier(sequence *atomicSequence) *atomicBarrier {
	return &atomicBarrier{sequence: sequence}
}
func (this *atomicBarrier) Load(_ int64) int64 { return this.sequence.Load() }

const defaultSequenceValue = -1
