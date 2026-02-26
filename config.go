package disruptor

import (
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

	if config.SequencerCount > 1 {
		return newMultiWriterDisruptor(config)
	}

	return newSingleWriterDisruptor(config)
}
func newSingleWriterDisruptor(config configuration) (*defaultDisruptor, error) {
	committedSequence := newSequence()
	listener, handledBarrier := config.newListeners(newAtomicBarrier(committedSequence))
	sequencer := newSequencer(config.BufferCapacity, committedSequence, handledBarrier, config.WaitStrategy)
	return &defaultDisruptor{ListenCloser: listener, sequencers: []Sequencer{sequencer}}, nil
}
func newMultiWriterDisruptor(config configuration) (*defaultDisruptor, error) {
	committed, shift := newCommittedBuffer(uint32(config.BufferCapacity))
	upper := newSequence()
	writeBarrier := newMultiSequencerBarrier(upper, committed, shift)
	listener, handledBarrier := config.newListeners(writeBarrier)

	sequencers := make([]Sequencer, config.SequencerCount)
	for i := range sequencers {
		sequencers[i] = newMultiSequencer(upper, committed, shift, handledBarrier, config.WaitStrategy)
	}

	return &defaultDisruptor{ListenCloser: listener, sequencers: sequencers}, nil
}
func (this configuration) newListeners(writeBarrier sequenceBarrier) (listener ListenCloser, handledBarrier sequenceBarrier) {
	handledBarrier = writeBarrier
	var listeners []ListenCloser

	for _, handlers := range this.HandlerGroups {
		group := make([]ListenCloser, 0, len(handlers))
		sequences := make([]*atomicSequence, 0, len(handlers))
		for _, handler := range handlers {
			currentSequence := newSequence()
			sequences = append(sequences, currentSequence)
			group = append(group, newListener(currentSequence, writeBarrier, handledBarrier, this.WaitStrategy, handler))
		}
		handledBarrier = newCompositeBarrier(sequences...) // next batch cannot handle beyond the sequences the current batch have handled.
		listeners = append(listeners, newCompositeListener(group))
	}

	return newCompositeListener(listeners), handledBarrier
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type configuration struct {
	BufferCapacity int64
	SequencerCount uint8
	WaitStrategy   WaitStrategy
	HandlerGroups  [][]Handler
}

func (singleton) BufferCapacity(value uint32) option {
	return func(this *configuration) { this.BufferCapacity = int64(value) }
}

// SequencerCount indicates the number of Sequencer instances to build. Each Sequencer instance should be attached to
// a message writer or producer and should not be shared among writers/producers without explicit synchronization.
func (singleton) SequencerCount(value uint8) option {
	return func(this *configuration) { this.SequencerCount = value }
}
func (singleton) WaitStrategy(value WaitStrategy) option {
	return func(this *configuration) { this.WaitStrategy = value }
}

// NewHandlerGroup defines a set of one or more Handler instances, each of which runs in its own goroutine, and which
// gate together. That is, each group does not allow a subsequent group of Handlers to operate on the underlying ring
// buffer until the current all Handlers within the current group have completed all operations.
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
		Options.SequencerCount(1),
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func newSequence() *atomicSequence {
	this := &atomicSequence{}
	this.Store(defaultSequenceValue)
	return this
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type defaultDisruptor struct {
	ListenCloser
	sequencers []Sequencer
}

func (this *defaultDisruptor) Sequencers() []Sequencer { return this.sequencers }
