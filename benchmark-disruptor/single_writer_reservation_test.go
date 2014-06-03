package benchmarks

import (
	"testing"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkDisruptorWriterReserveOne(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	reader := disruptor.NewReader(read, written, written, SampleConsumer{&ringBuffer})
	writer := disruptor.NewWriter(written, read, RingBufferSize)
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
func BenchmarkDisruptorWriterReserveMany(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	reader := disruptor.NewReader(read, written, written, SampleConsumer{&ringBuffer})
	writer := disruptor.NewWriter(written, read, RingBufferSize)
	reader.Start()

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	sequence := disruptor.InitialSequenceValue
	for sequence < iterations {
		sequence = writer.Reserve(ReserveMany)

		for i := sequence - ReserveManyDelta; i <= sequence; i++ {
			ringBuffer[i&RingBufferMask] = i
		}

		writer.Commit(sequence, sequence)
	}

	reader.Stop()
}
