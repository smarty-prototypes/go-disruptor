package disruptor

import (
	"errors"
	"sync/atomic"
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

	committedSequence := newSequence()
	listener, handledBarrier := config.newListeners(committedSequence)
	writer := newWriter(committedSequence, handledBarrier, config.BufferCapacity)

	return &defaultDisruptor{
		ListenCloser: listener,
		writers:      []Writer{writer}, // TODO: multi-writer
	}, nil
}
func (this configuration) newListeners(writeBarrier sequenceBarrier) (listener ListenCloser, handledBarrier sequenceBarrier) {
	handledBarrier = writeBarrier
	var listeners []ListenCloser

	for _, handlers := range this.HandlerGroups {
		group := make([]ListenCloser, 0, len(handlers))
		sequences := make([]*atomic.Int64, 0, len(handlers))
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
	WriterCount    uint8
	WaitStrategy   WaitStrategy
	HandlerGroups  [][]Handler
}

func (singleton) BufferCapacity(value uint32) option {
	return func(this *configuration) { this.BufferCapacity = int64(value) }
}
func (singleton) WriterCount(value uint8) option {
	return func(this *configuration) { this.WriterCount = value }
}
func (singleton) WaitStrategy(value WaitStrategy) option {
	return func(this *configuration) { this.WaitStrategy = value }
}
func (singleton) NewHandlerGroup(values ...Handler) option {
	// Each handler will be run on a separate goroutine. Each group of handlers will run as a unit meaning that a
	// subsequence group of handlers will not execute until the current group of handlers has completed successfully.
	return func(this *configuration) {
		filtered := make([]Handler, 0, len(values))
		for _, value := range values {
			if value != nil {
				filtered = append(filtered, value)
			}
		}

		if len(filtered) > 0 {
			this.HandlerGroups = append(this.HandlerGroups, values)
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
		Options.WriterCount(1),
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

func newSequence() *atomic.Int64 { return newAtomicInt64(defaultSequenceValue) }
func newAtomicInt64(initialState int64) *atomic.Int64 {
	this := &atomic.Int64{}
	this.Store(initialState)
	return this
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type defaultDisruptor struct {
	ListenCloser
	writers []Writer
}

func (this *defaultDisruptor) Writers() []Writer { return this.writers }
