package benchmarks

import (
	"runtime"
	"testing"
	"time"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkSharedWriterReserveOneContendedWrite(b *testing.B) {
	defer time.Sleep(DisruptorCleanup)
	runtime.GOMAXPROCS(3)
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
	defer time.Sleep(DisruptorCleanup)
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(1)

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
		current := disruptor.InitialSequenceValue
		for current < iterations {
			current = writer.Reserve(ReserveMany)

			for i := current - ReserveMany; i <= current; i++ {
				ringBuffer[i&RingBufferMask] = i
			}

			writer.Commit(current-ReserveMany, current)
		}
	}()
	current := disruptor.InitialSequenceValue
	for current < iterations {
		current = writer.Reserve(ReserveMany)

		for i := current - ReserveMany; i <= current; i++ {
			ringBuffer[i&RingBufferMask] = i
		}

		writer.Commit(current-ReserveMany, current)
	}

	b.StopTimer()
}

type noopConsumer struct{}

func (this noopConsumer) Consume(lower, upper int64) {}
