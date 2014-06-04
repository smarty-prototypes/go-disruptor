package benchmarks

import (
	"runtime"
	"testing"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkWriterReserveOneMultipleReaders(b *testing.B) {
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(2)
	ringBuffer := [RingBufferSize]int64{}
	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(SampleConsumer{&ringBuffer}, SampleConsumer{&ringBuffer}).
		Build()
	controller.Start()
	defer controller.Stop()
	writer := controller.Writer()

	iterations := int64(b.N)
	sequence := disruptor.InitialSequenceValue

	b.ReportAllocs()
	b.ResetTimer()

	for sequence < iterations {
		sequence = writer.Reserve(ReserveOne)
		ringBuffer[sequence&RingBufferMask] = sequence
		writer.Commit(sequence, sequence)
	}

	b.StopTimer()
}

func BenchmarkWriterReserveManyMultipleReaders(b *testing.B) {
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(2)
	ringBuffer := [RingBufferSize]int64{}
	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(SampleConsumer{&ringBuffer}, SampleConsumer{&ringBuffer}).
		Build()
	controller.Start()
	defer controller.Stop()
	writer := controller.Writer()

	iterations := int64(b.N)
	sequence := disruptor.InitialSequenceValue

	b.ReportAllocs()
	b.ResetTimer()

	for sequence < iterations {
		sequence = writer.Reserve(ReserveMany)

		for i := sequence - ReserveManyDelta; i <= sequence; i++ {
			ringBuffer[i&RingBufferMask] = i
		}

		writer.Commit(sequence, sequence)
	}

	b.StopTimer()
}
