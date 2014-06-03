package benchmarks

import (
	"testing"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkSharedWriterReserveOne(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	shared := disruptor.NewSharedWriterBarrier(written, RingBufferSize)
	reader := disruptor.NewReader(read, written, shared, SampleConsumer{&ringBuffer})
	writer := disruptor.NewSharedWriter(shared, read)
	reader.Start()

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	sequence := disruptor.InitialSequenceValue
	for sequence < iterations {
		sequence = writer.Reserve(ReserveOne)
		ringBuffer[sequence&RingBufferMask] = sequence
		writer.Commit(sequence, sequence)
	}

	reader.Stop()
}

func BenchmarkSharedWriterReserveMany(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	shared := disruptor.NewSharedWriterBarrier(written, RingBufferSize)
	reader := disruptor.NewReader(read, written, shared, SampleConsumer{&ringBuffer})
	writer := disruptor.NewSharedWriter(shared, read)
	reader.Start()

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

	reader.Stop()
}
