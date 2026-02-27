package disruptor

import (
	"log"
	"sync"
	"testing"
)

func BenchmarkChannelBlockingOneGoroutine(b *testing.B) {
	benchmarkBlocking(b, 1)
}
func BenchmarkChannelBlockingTwoGoroutines(b *testing.B) {
	benchmarkBlocking(b, 2)
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
	benchmarkNonBlocking(b, 2)
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

	for i := int64(0); i < maxReads; {
		select {
		case msg := <-channel:
			if writers == 1 && msg != i {
				panic("out of sequence")
			}
			i++
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

func BenchmarkSingleWriterReserve(b *testing.B) {
	read, written := newSequence(), newSequence()
	writer := newSequencer(1024, written, newAtomicBarrier(read), defaultWaitStrategy{})
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence := writer.Reserve(1)
		read.Store(sequence)
	}
}
func BenchmarkSingleWriterNextWrapPoint(b *testing.B) {
	read, written := newSequence(), newSequence()
	writer := newSequencer(1024*16, written, newAtomicBarrier(read), defaultWaitStrategy{})
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence := writer.Reserve(1)
		read.Store(sequence)
	}
}
func BenchmarkSingleWriterCommit(b *testing.B) {
	writer := newSequencer(1024, newSequence(), nil, defaultWaitStrategy{})
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		writer.Commit(i, i)
	}
}

func BenchmarkSingleWriterReserveOneSingleConsumer(b *testing.B) {
	benchmarkSequencerReservations(b, reserveOne, simpleHandler{})
}
func BenchmarkSingleWriterReserveManySingleConsumer(b *testing.B) {
	benchmarkSequencerReservations(b, reserveMany, simpleHandler{})
}
func BenchmarkSingleWriterReserveOneMultipleConsumers(b *testing.B) {
	benchmarkSequencerReservations(b, reserveOne, simpleHandler{}, simpleHandler{})
}
func BenchmarkSingleWriterReserveManyMultipleConsumers(b *testing.B) {
	benchmarkSequencerReservations(b, reserveMany, simpleHandler{}, simpleHandler{})
}
func benchmarkSequencerReservations(b *testing.B, count uint32, consumers ...Handler) {
	iterations := int64(b.N)

	simpleDisruptor := newSimpleDisruptor(false, consumers...)

	go func() {
		b.ReportAllocs()
		b.ResetTimer()

		var sequence int64 = -1
		slots := int64(count)
		for sequence < iterations {
			sequence = simpleDisruptor.Reserve(count)
			for i := sequence - (slots - 1); i <= sequence; i++ {
				ringBuffer[i&ringBufferMask] = i
			}
			simpleDisruptor.Commit(sequence-(slots-1), sequence)
		}

		_ = simpleDisruptor.Close()
	}()

	simpleDisruptor.Listen()
}

func BenchmarkSharedWriterReserveOneSingleConsumer(b *testing.B) {
	benchmarkSharedSequencerReservations(b, reserveOne, 2, simpleHandler{})
}
func BenchmarkSharedWriterReserveManySingleConsumer(b *testing.B) {
	benchmarkSharedSequencerReservations(b, reserveMany, 2, simpleHandler{})
}

func BenchmarkSharedWriterReserveOneMultipleConsumers(b *testing.B) {
	benchmarkSharedSequencerReservations(b, reserveOne, 2, simpleHandler{}, simpleHandler{})
}

func BenchmarkSharedWriterReserveManyMultipleConsumers(b *testing.B) {
	benchmarkSharedSequencerReservations(b, reserveMany, 2, simpleHandler{}, simpleHandler{})
}

func BenchmarkSharedWriterReserveManyMultipleConsumers_ThreeWriters(b *testing.B) {
	benchmarkSharedSequencerReservations(b, reserveMany, 3, simpleHandler{}, simpleHandler{})
}
func benchmarkSharedSequencerReservations(b *testing.B, count uint32, writerCount int64, consumers ...Handler) {
	iterations := int64(b.N)

	sharedDisruptor := newSimpleDisruptor(true, consumers...)

	go func() {
		b.ReportAllocs()
		b.ResetTimer()

		var waiter sync.WaitGroup
		waiter.Add(int(writerCount))

		slots := int64(count)
		for writerIndex := int64(0); writerIndex < writerCount; writerIndex++ {
			go func() {
				defer waiter.Done()
				var sequence int64 = -1
				for sequence < iterations {
					sequence = sharedDisruptor.Reserve(count)
					for i := sequence - (slots - 1); i <= sequence; i++ {
						ringBuffer[i&ringBufferMask] = i
					}
					sharedDisruptor.Commit(sequence-(slots-1), sequence)
				}
			}()
		}
		waiter.Wait()
		_ = sharedDisruptor.Close()
	}()

	sharedDisruptor.Listen()
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
func newSimpleDisruptor(shared bool, consumers ...Handler) Disruptor {
	this, _ := New(
		Options.BufferCapacity(ringBufferSize),
		Options.SingleWriter(!shared),
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
