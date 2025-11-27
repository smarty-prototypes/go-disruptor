package disruptor

import (
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
		newCursor(), newCursor(), newCursor(), newCursor())

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		barrier.Load()
	}
}

func BenchmarkCursorStore(b *testing.B) {
	iterations := int64(b.N)
	sequence := newCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence.Store(i)
	}
}
func BenchmarkCursorLoad(b *testing.B) {
	iterations := int64(b.N)
	sequence := newCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		_ = sequence.Load()
	}
}
func BenchmarkCursorLoadAsBarrier(b *testing.B) {
	var barrier sequenceBarrier = newCursor()
	iterations := int64(b.N)

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		_ = barrier.Load()
	}
}

func BenchmarkWriterReserve(b *testing.B) {
	read, written := newCursor(), newCursor()
	writer := newWriter(written, read, 1024)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence := writer.Reserve(1)
		read.Store(sequence)
	}
}
func BenchmarkWriterNextWrapPoint(b *testing.B) {
	read, written := newCursor(), newCursor()
	writer := newWriter(written, read, 1024)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	read.Store(math.MaxInt64)
	for i := int64(0); i < iterations; i++ {
		writer.Reserve(1)
	}
}
func BenchmarkWriterCommit(b *testing.B) {
	writer := newWriter(newCursor(), nil, 1024)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		writer.Commit(i, i)
	}
}

func BenchmarkWriterReserveOneSingleConsumer(b *testing.B) {
	benchmarkSequencerReservations(b, ReserveOne, SampleConsumer{})
}
func BenchmarkWriterReserveManySingleConsumer(b *testing.B) {
	benchmarkSequencerReservations(b, ReserveMany, SampleConsumer{})
}
func BenchmarkWriterReserveOneMultipleConsumers(b *testing.B) {
	benchmarkSequencerReservations(b, ReserveOne, SampleConsumer{}, SampleConsumer{})
}
func BenchmarkWriterReserveManyMultipleConsumers(b *testing.B) {
	benchmarkSequencerReservations(b, ReserveMany, SampleConsumer{}, SampleConsumer{})
}
func benchmarkSequencerReservations(b *testing.B, count int64, consumers ...Handler) {
	iterations := int64(b.N)

	myDisruptor := build(consumers...)

	go func() {
		b.ReportAllocs()
		b.ResetTimer()

		var sequence int64 = -1
		for sequence < iterations {
			sequence = myDisruptor.Reserve(count)
			for i := sequence - (count - 1); i <= sequence; i++ {
				ringBuffer[i&RingBufferMask] = i
			}
			myDisruptor.Commit(sequence, sequence)
		}

		_ = myDisruptor.Close()
	}()

	myDisruptor.Listen()
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SampleConsumer struct{}

func (this SampleConsumer) Handle(lower, upper int64) {
	var message int64
	for lower <= upper {
		message = ringBuffer[lower&RingBufferMask]
		if message != lower {
			log.Panicf("race condition: Sequence: %d, Message: %d", lower, message)
		}
		lower++
	}
}

func build(consumers ...Handler) *Disruptor {
	this, _ := New(
		Options.Capacity(RingBufferSize),
		Options.ConsumerGroup(consumers...))
	return this
}

const (
	RingBufferSize = 1024 * 64
	RingBufferMask = RingBufferSize - 1
	ReserveOne     = 1
	ReserveMany    = 16
)

var ringBuffer = [RingBufferSize]int64{}
