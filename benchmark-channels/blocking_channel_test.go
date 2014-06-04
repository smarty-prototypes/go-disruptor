package benchmarks

import (
	"runtime"
	"testing"
)

func BenchmarkBlockingOneGoroutine(b *testing.B) {
	runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(1)
	benchmarkBlocking(b, 1)
}

func BenchmarkBlockingTwoGoroutines(b *testing.B) {
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	benchmarkBlocking(b, 1)
}

func BenchmarkBlockingThreeGoroutinesWithContendedWrite(b *testing.B) {
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(1)
	benchmarkBlocking(b, 2)
}

func benchmarkBlocking(b *testing.B, writers int64) {
	iterations := int64(b.N)
	channel := make(chan int64, 1024*16)

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
			panic("Out of sequence")
		}
	}

	b.StopTimer()
}
