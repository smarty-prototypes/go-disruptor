package benchmarks

import (
	"runtime"
	"testing"
	"time"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkWriterReserveOne(b *testing.B) {
	defer time.Sleep(DisruptorCleanup)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)

	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(SampleConsumer{}).
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
func BenchmarkWriterReserveMany(b *testing.B) {
	defer time.Sleep(DisruptorCleanup)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)

	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(SampleConsumer{}).
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
