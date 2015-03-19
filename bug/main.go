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
		previous, current := disruptor.InitialSequenceValue, disruptor.InitialSequenceValue
		for current < iterations {
			current = writer.Reserve("A", ReserveMany)

			for i := previous + 1; i <= current; i++ {
				ringBuffer[i&RingBufferMask] = i
			}

			writer.Commit("A", previous+1, current)
			previous = current
		}
	}()
	{
		previous, current := disruptor.InitialSequenceValue, disruptor.InitialSequenceValue
		for current < iterations {
			current = writer.Reserve("B", ReserveMany)

			for i := previous + 1; i <= current; i++ {
				ringBuffer[i&RingBufferMask] = i
			}

			writer.Commit("B", previous+1, current)
			previous = current
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
