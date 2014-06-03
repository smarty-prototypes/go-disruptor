package benchmarks

import (
	"testing"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkWriterReserveOneMultipleReaders(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	written, read1, read2 := disruptor.NewCursor(), disruptor.NewCursor(), disruptor.NewCursor()
	reader1 := disruptor.NewReader(read1, written, written, SampleConsumer{&ringBuffer})
	reader2 := disruptor.NewReader(read2, written, written, SampleConsumer{&ringBuffer})
	barrier := disruptor.NewCompositeBarrier(read1, read2)
	writer := disruptor.NewWriter(written, barrier, RingBufferSize)

	reader1.Start()
	reader2.Start()

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	sequence := disruptor.InitialSequenceValue
	for sequence < iterations {
		sequence = writer.Reserve(ReserveOne)
		ringBuffer[sequence&RingBufferMask] = sequence
		writer.Commit(sequence, sequence)
	}

	reader1.Stop()
	reader2.Stop()
}
func BenchmarkWriterReserveManyMultipleReaders(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	written, read1, read2 := disruptor.NewCursor(), disruptor.NewCursor(), disruptor.NewCursor()
	reader1 := disruptor.NewReader(read1, written, written, SampleConsumer{&ringBuffer})
	reader2 := disruptor.NewReader(read2, written, written, SampleConsumer{&ringBuffer})
	barrier := disruptor.NewCompositeBarrier(read1, read2)
	writer := disruptor.NewWriter(written, barrier, RingBufferSize)

	reader1.Start()
	reader2.Start()

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

	reader1.Stop()
	reader2.Stop()
}
