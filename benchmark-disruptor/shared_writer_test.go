package benchmarks

import (
	"runtime"
	"testing"
	"time"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkSharedWriterReserveOne(b *testing.B) {
	defer time.Sleep(DisruptorCleanup)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)

	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(SampleConsumer{}).
		BuildShared()
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

func BenchmarkSharedWriterReserveMany(b *testing.B) {
	defer time.Sleep(DisruptorCleanup)
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)

	controller := disruptor.
		Configure(RingBufferSize).
		WithConsumerGroup(SampleConsumer{}).
		BuildShared()
	controller.Start()
	defer controller.Stop()
	writer := controller.Writer()

	iterations := int64(b.N)

	b.ReportAllocs()
	b.ResetTimer()

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
