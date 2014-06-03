package benchmarks

import (
	"fmt"
	"testing"

	"github.com/smartystreets/go-disruptor"
)

func BenchmarkDisruptorWriterReserveSingle(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	reader := disruptor.NewReader(read, written, written, Consumer{&ringBuffer})
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
func BenchmarkDisruptorWriterReserveMultiple(b *testing.B) {
	ringBuffer := [RingBufferSize]int64{}
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	reader := disruptor.NewReader(read, written, written, Consumer{&ringBuffer})
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

type Consumer struct {
	ringBuffer *[RingBufferSize]int64
}

func (this Consumer) Consume(lower, upper int64) {
	for lower <= upper {
		message := this.ringBuffer[lower&RingBufferMask]
		if message != lower {
			panic(fmt.Sprintf("\nRace condition %d %d\n", lower, message))
		}
		lower++
	}
}
