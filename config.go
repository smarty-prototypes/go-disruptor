package disruptor

import (
	"context"
	"errors"
	"runtime"
	"time"
)

func New(options ...option) (Disruptor, error) {
	config := configuration{}
	Options.apply(options...)(&config)

	if config.BufferCapacity <= 0 {
		return nil, errors.New("buffer capacity must be at least 1")
	} else if config.BufferCapacity&(config.BufferCapacity-1) != 0 {
		return nil, errors.New("the buffer capacity must be a power of two, e.g. 2, 4, 8, 16")
	} else if len(config.HandlerGroups) == 0 {
		return nil, errors.New("no handlers have been provided")
	}

	upperSequence := newSequence()
	if config.SingleWriter {
		listener, handledBarrier := config.newListeners(newAtomicBarrier(upperSequence))
		sequencer := newSequencer(config.BufferCapacity, upperSequence, handledBarrier, config.WaitStrategy)
		return &defaultDisruptor{ListenCloser: listener, Sequencer: sequencer}, nil
	}

	sequencer := newSharedSequencer(uint32(config.BufferCapacity), upperSequence, config.WaitStrategy)
	listener, handledBarrier := config.newListeners(sequencer)
	sequencer.consumerBarrier = handledBarrier
	return &defaultDisruptor{ListenCloser: listener, Sequencer: sequencer}, nil
}
func (this configuration) newListeners(committedBarrier sequenceBarrier) (listener ListenCloser, handledBarrier sequenceBarrier) {
	handledBarrier = committedBarrier

	totalSequences := 0
	for _, handlers := range this.HandlerGroups {
		totalSequences += len(handlers)
	}
	allSequences := newSequences(totalSequences)

	listeners := make([]ListenCloser, len(this.HandlerGroups))
	offset := 0
	for groupIndex, handlers := range this.HandlerGroups {
		sequences := make([]*atomicSequence, len(handlers))
		group := make([]ListenCloser, len(handlers))
		for handlerIndex, handler := range handlers {
			sequences[handlerIndex] = allSequences[offset+handlerIndex]
			group[handlerIndex] = newListener(sequences[handlerIndex], committedBarrier, handledBarrier, this.WaitStrategy, handler)
		}
		handledBarrier = newCompositeBarrier(sequences...) // next group cannot handle beyond the sequences the current group have handled.
		listeners[groupIndex] = newCompositeListener(group)
		offset += len(handlers)
	}

	return newCompositeListener(listeners), handledBarrier
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type configuration struct {
	BufferCapacity uint32
	SingleWriter   bool
	WaitStrategy   WaitStrategy
	HandlerGroups  [][]Handler
}

// BufferCapacity sets the number of slots in the ring buffer. Must be a power of 2. Default: 1024.
func (singleton) BufferCapacity(value uint32) option {
	return func(this *configuration) { this.BufferCapacity = value }
}

// SingleWriter configures whether the Disruptor uses a single-writer Sequencer (true) or a multi-writer shared
// Sequencer (false). The single-writer Sequencer is faster but must not be used concurrently from multiple
// goroutines. Default: true.
func (singleton) SingleWriter(value bool) option {
	return func(this *configuration) { this.SingleWriter = value }
}

// WaitStrategy sets the backpressure strategy used by both producers and consumers. Default: defaultWaitStrategy.
func (singleton) WaitStrategy(value WaitStrategy) option {
	return func(this *configuration) { this.WaitStrategy = value }
}

// NewHandlerGroup defines a set of one or more Handler instances, each of which runs in its own goroutine, and which
// gate together. That is, each group does not allow a subsequent group of Handlers to operate on the underlying ring
// buffer until all Handlers within the current group have completed all operations.
func (singleton) NewHandlerGroup(values ...Handler) option {
	return func(this *configuration) {
		filtered := make([]Handler, 0, len(values))
		for _, value := range values {
			if value != nil {
				filtered = append(filtered, value)
			}
		}

		if len(filtered) > 0 {
			this.HandlerGroups = append(this.HandlerGroups, filtered)
		}
	}
}
func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		for _, item := range Options.defaults(options...) {
			item(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	return append([]option{
		Options.BufferCapacity(1024),
		Options.SingleWriter(true),
		Options.WaitStrategy(defaultWaitStrategy{}),
	}, options...)
}

type singleton struct{}
type option func(*configuration)

var Options singleton

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type defaultWaitStrategy struct{}

func (this defaultWaitStrategy) Gate(int64) { time.Sleep(time.Nanosecond) }
func (this defaultWaitStrategy) Idle(int64) { time.Sleep(time.Millisecond) }
func (this defaultWaitStrategy) Reserve()   { runtime.Gosched() } // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
func (this defaultWaitStrategy) TryReserve(ctx context.Context) error {
	runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
	return ctx.Err()
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func newSequence() *atomicSequence {
	this := &atomicSequence{}
	this.Store(defaultSequenceValue)
	return this
}

// newSequences allocates a slice of *atomicSequence in a contiguous space in memory
func newSequences(count int) []*atomicSequence {
	backing := make([]atomicSequence, count)
	sequences := make([]*atomicSequence, count)
	for i := range sequences {
		backing[i].Store(defaultSequenceValue)
		sequences[i] = &backing[i]
	}
	return sequences
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type defaultDisruptor struct {
	ListenCloser
	Sequencer
}
