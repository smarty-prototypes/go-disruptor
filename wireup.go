package disruptor

import (
	"errors"
	"sync/atomic"
)

type configuration struct {
	WaitStrategy   WaitStrategy
	Capacity       int64
	ListenerGroups [][]Handler
}

func New(options ...option) (Disruptor, error) {
	config := configuration{}
	Options.apply(options...)(&config)

	if config.Capacity <= 0 {
		return nil, errCapacityTooSmall
	}

	if config.Capacity&(config.Capacity-1) != 0 {
		return nil, errCapacityPowerOfTwo
	}

	if len(config.ListenerGroups) == 0 {
		return nil, errNoListeners
	}

	for _, listenerGroup := range config.ListenerGroups {
		if len(listenerGroup) == 0 {
			return nil, errEmptyListenerGroup
		}

		for _, consumer := range listenerGroup {
			if consumer == nil {
				return nil, errEmptyListener
			}
		}
	}

	var writerSequence = newCursor()
	listeners, readBarrier := config.newListeners(writerSequence)

	return struct {
		Writer
		ListenCloser
	}{
		Writer:       newWriter(writerSequence, readBarrier, config.Capacity),
		ListenCloser: compositeListener(listeners),
	}, nil
}
func (this configuration) newListeners(writerSequence *atomic.Int64) (listeners []ListenCloser, upstream sequenceBarrier) {
	upstream = writerSequence

	for _, listenerGroup := range this.ListenerGroups {
		var consumerGroupSequences []*atomic.Int64

		for _, consumer := range listenerGroup {
			currentSequence := newCursor()
			listeners = append(listeners, newListener(currentSequence, writerSequence, upstream, this.WaitStrategy, consumer))
			consumerGroupSequences = append(consumerGroupSequences, currentSequence)
		}

		upstream = newCompositeBarrier(consumerGroupSequences...)
	}

	return listeners, upstream
}

func (singleton) WaitStrategy(value WaitStrategy) option {
	return func(this *configuration) { this.WaitStrategy = value }
}
func (singleton) Capacity(value int64) option {
	return func(this *configuration) { this.Capacity = value }
}
func (singleton) NewListenerGroup(values ...Handler) option {
	return func(this *configuration) { this.ListenerGroups = append(this.ListenerGroups, values) }
}

func (singleton) apply(options ...option) option {
	return func(this *configuration) {
		for _, item := range Options.defaults(options...) {
			item(this)
		}
	}
}
func (singleton) defaults(options ...option) []option {
	const defaultCapacity = 1024
	var waitStrategy = defaultWaitStrategy{}

	return append([]option{
		Options.Capacity(defaultCapacity),
		Options.WaitStrategy(waitStrategy),
	}, options...)
}

type singleton struct{}
type option func(*configuration)

var Options singleton

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

var (
	errCapacityTooSmall   = errors.New("the capacity must be at least 1")
	errCapacityPowerOfTwo = errors.New("the capacity be a power of two, e.g. 2, 4, 8, 16")
	errNoListeners        = errors.New("no consumers have been provided")
	errEmptyListenerGroup = errors.New("empty listener group")
	errEmptyListener      = errors.New("an empty listener was specified in the listener group")
)
