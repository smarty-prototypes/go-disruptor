package disruptor

import (
	"errors"
	"sync/atomic"
)

type configuration struct {
	WaitStrategy   WaitStrategy
	Capacity       int64
	ConsumerGroups [][]Handler
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

	if len(config.ConsumerGroups) == 0 {
		return nil, errMissingConsumers
	}

	for _, consumerGroup := range config.ConsumerGroups {
		if len(consumerGroup) == 0 {
			return nil, errMissingConsumersInGroup
		}

		for _, consumer := range consumerGroup {
			if consumer == nil {
				return nil, errEmptyConsumer
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

	for _, consumerGroup := range this.ConsumerGroups {
		var consumerGroupSequences []*atomic.Int64

		for _, consumer := range consumerGroup {
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
func (singleton) AddConsumerGroup(values ...Handler) option {
	return func(this *configuration) { this.ConsumerGroups = append(this.ConsumerGroups, values) }
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
	errCapacityTooSmall        = errors.New("the capacity must be at least 1")
	errCapacityPowerOfTwo      = errors.New("the capacity be a power of two, e.g. 2, 4, 8, 16")
	errMissingConsumers        = errors.New("no consumers have been provided")
	errMissingConsumersInGroup = errors.New("the consumer group does not have any consumers")
	errEmptyConsumer           = errors.New("an empty consumer was specified in the consumer group")
)
