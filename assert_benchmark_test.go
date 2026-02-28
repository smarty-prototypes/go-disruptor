package disruptor

import (
	"sync"
	"testing"
	"time"
)

// Key:
// SP   = single-producer
// MP   = multi-producer
// SC   = single-consumer
// MC   = multi-producer
// R[N] = reserve N number of slots

func BenchmarkChannel(b *testing.B) {
	b.Run("SPSC/Blocking", func(b *testing.B) { benchmarkChannelBlocking(b, 1) })
	b.Run("MPSC/Blocking", func(b *testing.B) { benchmarkChannelBlocking(b, 2) })
	b.Run("SPSC/NonBlocking", func(b *testing.B) { benchmarkChannelNonBlocking(b, 1) })
	b.Run("MPSC/NonBlocking", func(b *testing.B) { benchmarkChannelNonBlocking(b, 2) })
}
func benchmarkChannelBlocking(b *testing.B, writers int64) {
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
func benchmarkChannelNonBlocking(b *testing.B, writers int64) {
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

func BenchmarkSequence(b *testing.B) {
	b.Run("Store", func(b *testing.B) {
		sequence := newSequence()
		b.ReportAllocs()
		b.ResetTimer()
		for i := int64(0); i < int64(b.N); i++ {
			sequence.Store(i)
		}
	})
	b.Run("Load", func(b *testing.B) {
		sequence := newSequence()
		b.ReportAllocs()
		b.ResetTimer()
		for i := int64(0); i < int64(b.N); i++ {
			_ = sequence.Load()
		}
	})
	b.Run("LoadAsBarrier", func(b *testing.B) {
		var barrier sequenceBarrier = newAtomicBarrier(newSequence())
		b.ReportAllocs()
		b.ResetTimer()
		for i := int64(0); i < int64(b.N); i++ {
			_ = barrier.Load(0)
		}
	})
	b.Run("CompositeBarrier", func(b *testing.B) {
		barrier := newCompositeBarrier(newSequence(), newSequence(), newSequence(), newSequence())
		b.ReportAllocs()
		b.ResetTimer()
		for i := int64(0); i < int64(b.N); i++ {
			barrier.Load(0)
		}
	})
}
func BenchmarkSequencer(b *testing.B) {
	b.Run("Reserve", func(b *testing.B) {
		read, written := newSequence(), newSequence()
		writer := newSequencer(1024, written, newAtomicBarrier(read), defaultWaitStrategy{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := int64(0); i < int64(b.N); i++ {
			sequence := writer.Reserve(1)
			read.Store(sequence)
		}
	})
	b.Run("NextWrapPoint", func(b *testing.B) {
		read, written := newSequence(), newSequence()
		writer := newSequencer(1024*16, written, newAtomicBarrier(read), defaultWaitStrategy{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := int64(0); i < int64(b.N); i++ {
			sequence := writer.Reserve(1)
			read.Store(sequence)
		}
	})
	b.Run("Commit", func(b *testing.B) {
		writer := newSequencer(1024, newSequence(), nil, defaultWaitStrategy{})
		b.ReportAllocs()
		b.ResetTimer()
		for i := int64(0); i < int64(b.N); i++ {
			writer.Commit(i, i)
		}
	})
	b.Run("SP SC", func(b *testing.B) { benchmarkDisruptor(b, false, reserve1, 1, nopHandler{}) })
	b.Run("SP4SC", func(b *testing.B) { benchmarkDisruptor(b, false, reserve4, 1, nopHandler{}) })
	b.Run("SP MC", func(b *testing.B) { benchmarkDisruptor(b, false, reserve1, 1, nopHandler{}, nopHandler{}) })
	b.Run("SP4MC", func(b *testing.B) { benchmarkDisruptor(b, false, reserve4, 1, nopHandler{}, nopHandler{}) })
}
func BenchmarkSharedSequencer(b *testing.B) {
	b.Run("SP SC/R1", func(b *testing.B) { benchmarkDisruptor(b, true, reserve1, 1, nopHandler{}) })
	b.Run("MP SC/R1", func(b *testing.B) { benchmarkDisruptor(b, true, reserve1, 4, nopHandler{}) })
	b.Run("MP SC/R4", func(b *testing.B) { benchmarkDisruptor(b, true, reserve4, 4, nopHandler{}) })
	b.Run("MP MC/R1", func(b *testing.B) { benchmarkDisruptor(b, true, reserve1, 4, nopHandler{}, nopHandler{}) })
	b.Run("MP MC/R4", func(b *testing.B) { benchmarkDisruptor(b, true, reserve4, 4, nopHandler{}, nopHandler{}) })
}
func BenchmarkSharedOrderedSequencer(b *testing.B) {
	mk := func(b *testing.B, reserve uint32, writers int, consumers ...Handler) {
		benchmarkDisruptorWith(b, reserve, writers, consumers, Options.WriterContention(ContentionLow))
	}
	b.Run("SP SC/R1", func(b *testing.B) { mk(b, reserve1, 1, nopHandler{}) })
	b.Run("MP SC/R1", func(b *testing.B) { mk(b, reserve1, 2, nopHandler{}) })
	b.Run("MP SC/R4", func(b *testing.B) { mk(b, reserve4, 2, nopHandler{}) })
	b.Run("MP MC/R1", func(b *testing.B) { mk(b, reserve1, 2, nopHandler{}, nopHandler{}) })
	b.Run("MP MC/R4", func(b *testing.B) { mk(b, reserve4, 2, nopHandler{}, nopHandler{}) })
	b.Run("MP3MC/R4", func(b *testing.B) { mk(b, reserve4, 3, nopHandler{}, nopHandler{}) })
}
func benchmarkDisruptor(b *testing.B, shared bool, count uint32, writerCount int, consumers ...Handler) {
	benchmarkDisruptorWith(b, count, writerCount, consumers, Options.WriterContention(WriterContention(writerCount)))
}
func benchmarkDisruptorWith(b *testing.B, count uint32, writerCount int, consumers []Handler, opts ...option) {
	iterations := int64(b.N)
	slots := int64(count)
	offset := slots - 1

	opts = append([]option{Options.BufferCapacity(ringBufferSize), Options.NewHandlerGroup(consumers...)}, opts...)
	disruptor, _ := New(opts...)
	defer disruptor.Listen()

	go func() {
		var waiter sync.WaitGroup
		waiter.Add(writerCount)
		defer func() { waiter.Wait(); _ = disruptor.Close() }()
		time.Sleep(time.Millisecond * 100) // let the Listen goroutine have time to start
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < writerCount; i++ {
			go func() {
				defer waiter.Done()
				for sequence := int64(defaultSequenceValue); sequence < iterations; {
					sequence = disruptor.Reserve(count)
					disruptor.Commit(sequence-offset, sequence)
				}
			}()
		}
	}()
}

type nopHandler struct{}

func (this nopHandler) Handle(int64, int64) {}

const (
	ringBufferSize = 1 << 16 // 64K
	reserve1       = 1
	reserve4       = 16
)
