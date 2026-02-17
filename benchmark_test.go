package disruptor

import (
	"context"
	"log"
	"math"
	"testing"
)

func BenchmarkChannelBlockingOneGoroutine(b *testing.B) {
	benchmarkBlocking(b, 1)
}
func BenchmarkChannelBlockingTwoGoroutines(b *testing.B) {
	benchmarkBlocking(b, 1)
}
func BenchmarkChannelBlockingThreeGoroutinesWithContendedWrite(b *testing.B) {
	benchmarkBlocking(b, 2)
}
func benchmarkBlocking(b *testing.B, writers int64) {
	channel := make(chan int64, 1024*16)
	iterations := int64(b.N)

	b.ReportAllocs()
	b.ResetTimer()

	for x := int64(0); x < writers; x++ {
		go func() {
			for i := int64(0); i < iterations; i++ {
				channel <- i
			}
		}()
	}

	for i := int64(0); i < iterations*writers; i++ {
		msg := <-channel
		if writers == 1 && msg != i {
			panic("out of sequence")
		}
	}
}

func BenchmarkChannelNonBlockingOneGoroutine(b *testing.B) {
	benchmarkNonBlocking(b, 1)
}
func BenchmarkChannelNonBlockingTwoGoroutines(b *testing.B) {
	benchmarkNonBlocking(b, 1)
}
func BenchmarkChannelNonBlockingThreeGoroutinesWithContendedWrite(b *testing.B) {
	benchmarkNonBlocking(b, 2)
}
func benchmarkNonBlocking(b *testing.B, writers int64) {
	iterations := int64(b.N)
	maxReads := iterations * writers
	channel := make(chan int64, 1024*16)

	b.ReportAllocs()
	b.ResetTimer()

	for x := int64(0); x < writers; x++ {
		go func() {
			for i := int64(0); i < iterations; {
				select {
				case channel <- i:
					i++
				default:
					continue
				}
			}
		}()
	}

	for i := int64(0); i < maxReads; i++ {
		select {
		case msg := <-channel:
			if writers == 1 && msg != i {
				// panic("Out of sequence")
			}
		default:
			continue
		}
	}
}

func BenchmarkCompositeBarrierRead(b *testing.B) {
	iterations := int64(b.N)

	barrier := newCompositeBarrier(
		newSequence(), newSequence(), newSequence(), newSequence())

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		barrier.Load(0)
	}
}

func BenchmarkSequenceStore(b *testing.B) {
	iterations := int64(b.N)
	sequence := newSequence()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence.Store(i)
	}
}
func BenchmarkSequenceLoad(b *testing.B) {
	iterations := int64(b.N)
	sequence := newSequence()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		_ = sequence.Load()
	}
}
func BenchmarkSequenceLoadAsBarrier(b *testing.B) {
	var barrier sequenceBarrier = newAtomicBarrier(newSequence())
	iterations := int64(b.N)

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		_ = barrier.Load(0)
	}
}

func BenchmarkWriterReserve(b *testing.B) {
	read, written := newSequence(), newSequence()
	writer := newSequencer(written, newAtomicBarrier(read), 1024)
	ctx := context.Background()
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence := writer.Reserve(ctx, 1)
		read.Store(sequence)
	}
}
func BenchmarkWriterNextWrapPoint(b *testing.B) {
	read, written := newSequence(), newSequence()
	writer := newSequencer(written, newAtomicBarrier(read), 1024)
	ctx := context.Background()
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	read.Store(math.MaxInt64)
	for i := int64(0); i < iterations; i++ {
		writer.Reserve(ctx, 1)
	}
}
func BenchmarkWriterCommit(b *testing.B) {
	writer := newSequencer(newSequence(), nil, 1024)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		writer.Commit(i, i)
	}
}

func BenchmarkWriterReserveOneSingleConsumer(b *testing.B) {
	benchmarkSequencerReservations(b, reserveOne, simpleHandler{})
}
func BenchmarkWriterReserveManySingleConsumer(b *testing.B) {
	benchmarkSequencerReservations(b, reserveMany, simpleHandler{})
}
func BenchmarkWriterReserveOneMultipleConsumers(b *testing.B) {
	benchmarkSequencerReservations(b, reserveOne, simpleHandler{}, simpleHandler{})
}
func BenchmarkWriterReserveManyMultipleConsumers(b *testing.B) {
	benchmarkSequencerReservations(b, reserveMany, simpleHandler{}, simpleHandler{})
}
func benchmarkSequencerReservations(b *testing.B, count int64, consumers ...Handler) {
	iterations := int64(b.N)

	simpleDisruptor := newSimpleDisruptor(consumers...)
	writer := simpleDisruptor.Sequencers()[0]
	ctx := context.Background()

	go func() {
		b.ReportAllocs()
		b.ResetTimer()

		var sequence int64 = -1
		for sequence < iterations {
			sequence = writer.Reserve(ctx, count)
			for i := sequence - (count - 1); i <= sequence; i++ {
				ringBuffer[i&ringBufferMask] = i
			}
			writer.Commit(sequence, sequence)
		}

		_ = simpleDisruptor.Close()
	}()

	simpleDisruptor.Listen()
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type simpleHandler struct{}

func (this simpleHandler) Handle(lower, upper int64) {
	var message int64
	for lower <= upper {
		message = ringBuffer[lower&ringBufferMask]
		if message != lower {
			log.Panicf("race condition: Sequence: %d, Message: %d", lower, message)
		}
		lower++
	}
}
func newSimpleDisruptor(consumers ...Handler) Disruptor {
	this, _ := New(
		Options.BufferCapacity(ringBufferSize),
		Options.NewHandlerGroup(consumers...))
	return this
}

const (
	ringBufferSize = 1024 * 64
	ringBufferMask = ringBufferSize - 1
	reserveOne     = 1
	reserveMany    = 16
)

var ringBuffer = [ringBufferSize]int64{}
