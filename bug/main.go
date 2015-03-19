package main

import (
	"fmt"
	"runtime"
	"time"

	"github.com/smartystreets/go-disruptor"
)

func main() {
	runtime.GOMAXPROCS(3)

	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(noopConsumer{}).
		BuildShared()

	controller.Start()
	defer controller.Stop()
	writer := controller.Writer()

	go func() {
		current := disruptor.InitialSequenceValue
		for current < iterations {
			current = writer.Reserve("A", ReserveMany)

			for i := current - ReserveMany; i <= current; i++ {
				ringBuffer[i&RingBufferMask] = i
			}

			writer.Commit("A", current-ReserveMany, current)
		}
	}()
	{
		current := disruptor.InitialSequenceValue
		for current < iterations {
			current = writer.Reserve("B", ReserveMany)

			for i := current - ReserveMany; i <= current; i++ {
				ringBuffer[i&RingBufferMask] = i
			}

			writer.Commit("B", current-ReserveMany, current)
		}
	}
}

type noopConsumer struct{}

func (this noopConsumer) Consume(lower, upper int64) {
	fmt.Printf("Consumer consumed up to %d\n", upper)
}

const iterations = 1024 * 1024

const (
	RingBufferSize   = 64
	RingBufferMask   = RingBufferSize - 1
	ReserveOne       = 1
	ReserveMany      = 16
	ReserveManyDelta = ReserveMany - 1
	DisruptorCleanup = time.Millisecond * 10
)

var ringBuffer = [RingBufferSize]int64{}
