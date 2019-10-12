package benchmarks

import (
	"log"
	"math"
	"testing"

	"github.com/smartystreets-prototypes/go-disruptor"
)

func BenchmarkSequencerReserve(b *testing.B) {
	read, written := disruptor.NewCursor(), disruptor.NewCursor()
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
	read, written := disruptor.NewCursor(), disruptor.NewCursor()
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
	writer := disruptor.NewSequencer(disruptor.NewCursor(), nil, 1024)
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

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SampleConsumer struct{}

func (this SampleConsumer) Consume(lower, upper int64) {
	var message int64
	for lower <= upper {
		message = ringBuffer[lower&RingBufferMask]
		if message != lower {
			log.Panicf("race condition: Sequence: %d, Message: %d", lower, message)
		}
		lower++
	}
}

func build(consumers ...disruptor.Consumer) (disruptor.Sequencer, disruptor.ListenCloser) {
	return disruptor.New(
		disruptor.WithCapacity(RingBufferSize),
		disruptor.WithConsumerGroup(consumers...))
}

const (
	RingBufferSize = 1024 * 64
	RingBufferMask = RingBufferSize - 1
	ReserveOne     = 1
	ReserveMany    = 16
)

var ringBuffer = [RingBufferSize]int64{}
