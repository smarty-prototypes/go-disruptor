package benchmarks

import (
	"testing"

	"github.com/smartystreets-prototypes/go-disruptor"
)

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
	sequencer, listener := build(consumers...)

	go func() {
		b.ReportAllocs()
		b.ResetTimer()

		iterations := int64(b.N)
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
	b.StopTimer()
}
