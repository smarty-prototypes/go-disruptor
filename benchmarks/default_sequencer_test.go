package benchmarks

import (
	"math"
	"testing"

	"github.com/smartystreets-prototypes/go-disruptor"
)

func BenchmarkSequencerReserve(b *testing.B) {
	read, written := disruptor.NewSequence(), disruptor.NewSequence()
	writer := disruptor.NewSequencer(written, read, 1024)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence := writer.Reserve(1)
		read.Store(sequence)
	}
}
func BenchmarkSequencerNextWrapPoint(b *testing.B) {
	read, written := disruptor.NewSequence(), disruptor.NewSequence()
	writer := disruptor.NewSequencer(written, read, 1024)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	read.Store(math.MaxInt64)
	for i := int64(0); i < iterations; i++ {
		writer.Reserve(1)
	}
}
func BenchmarkSequencerCommit(b *testing.B) {
	writer := disruptor.NewSequencer(disruptor.NewSequence(), nil, 1024)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		writer.Commit(i, i)
	}
}

func BenchmarkSequencerReserveOneSingleConsumer(b *testing.B) {
	benchmarkSequencerReservations(b, ReserveOne, SampleConsumer{})
}
func BenchmarkSequencerReserveManySingleConsumer(b *testing.B) {
	benchmarkSequencerReservations(b, ReserveMany, SampleConsumer{})
}
func BenchmarkSequencerReserveOneMultipleConsumers(b *testing.B) {
	benchmarkSequencerReservations(b, ReserveOne, SampleConsumer{}, SampleConsumer{})
}
func BenchmarkSequencerReserveManyMultipleConsumers(b *testing.B) {
	benchmarkSequencerReservations(b, ReserveMany, SampleConsumer{}, SampleConsumer{})
}
func benchmarkSequencerReservations(b *testing.B, count int64, consumers ...disruptor.Consumer) {
	iterations := int64(b.N)

	sequencer, listener := build(consumers...)

	go func() {
		b.ReportAllocs()
		b.ResetTimer()

		var sequence int64 = -1
		for sequence < iterations {
			sequence = sequencer.Reserve(count)
			for i := sequence - (count - 1); i <= sequence; i++ {
				ringBuffer[i&RingBufferMask] = i
			}
			sequencer.Commit(sequence, sequence)
		}

		_ = listener.Close()
	}()

	listener.Listen()
}

