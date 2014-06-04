package benchmarks

import (
	"testing"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkWriterAwaitOne(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(SampleConsumer{&ringBuffer}).
		Build()
	controller.Start()
	defer controller.Stop()
	writer := controller.Writer()

	iterations := int64(b.N)
	sequence := disruptor.InitialSequenceValue

	b.ReportAllocs()
	b.ResetTimer()

	for sequence < iterations {
		sequence += ReserveOne
		writer.Await(sequence)
		ringBuffer[sequence&RingBufferMask] = sequence
		writer.Commit(sequence, sequence)
	}

	b.StopTimer()
}
func BenchmarkWriterAwaitMany(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(SampleConsumer{&ringBuffer}).
		Build()
	controller.Start()
	defer controller.Stop()
	writer := controller.Writer()

	iterations := int64(b.N)
	sequence := disruptor.InitialSequenceValue

	b.ReportAllocs()
	b.ResetTimer()

	for sequence < iterations {
		sequence += ReserveMany
		writer.Await(sequence)

		for i := sequence - ReserveManyDelta; i <= sequence; i++ {
			ringBuffer[i&RingBufferMask] = i
		}

		writer.Commit(sequence, sequence)
	}

	b.StopTimer()
}
