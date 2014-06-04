package benchmarks

import (
	"runtime"
	"testing"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkSharedWriterReserveOneContendedWrite(b *testing.B) {
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(2)

	ringBuffer := [RingBufferSize]int64{}
	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(SampleConsumer{&ringBuffer}).
		BuildShared()
	controller.Start()
	defer controller.Stop()
	writer := controller.Writer()

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	go func() {
		sequence := disruptor.InitialSequenceValue
		for sequence < iterations {
			sequence = writer.Reserve(ReserveOne)
			ringBuffer[sequence&RingBufferMask] = sequence
			writer.Commit(sequence, sequence)
		}
	}()

	sequence := disruptor.InitialSequenceValue
	for sequence < iterations {
		sequence = writer.Reserve(ReserveOne)
		ringBuffer[sequence&RingBufferMask] = sequence
		writer.Commit(sequence, sequence)
	}

	b.StopTimer()
}

func BenchmarkSharedWriterReserveManyContendedWrite(b *testing.B) {
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(2)

	ringBuffer := [RingBufferSize]int64{}
	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(noopConsumer{}).
		BuildShared()
	controller.Start()
	defer controller.Stop()
	writer := controller.Writer()

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	go func() {
		previous, current := disruptor.InitialSequenceValue, disruptor.InitialSequenceValue
		for current < iterations {
			current = writer.Reserve(ReserveMany)

			for i := previous + 1; i <= current; i++ {
				ringBuffer[i&RingBufferMask] = i
			}

			writer.Commit(previous+1, current)
			previous = current
		}
	}()
	previous, current := disruptor.InitialSequenceValue, disruptor.InitialSequenceValue
	for current < iterations {
		current = writer.Reserve(ReserveMany)

		for i := previous + 1; i <= current; i++ {
			ringBuffer[i&RingBufferMask] = i
		}

		writer.Commit(previous+1, current)
		previous = current
	}

	b.StopTimer()
}

type noopConsumer struct{}

func (this noopConsumer) Consume(lower, upper int64) {}
