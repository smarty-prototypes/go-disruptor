package benchmarks

import (
	"runtime"
	"testing"
)

func BenchmarkBlockingChannelOneGoroutine(b *testing.B) {
	benchmarkBlockingChannel(b)
}

func BenchmarkBlockingChannelTwoGoroutines(b *testing.B) {
	runtime.GOMAXPROCS(2)
	benchmarkBlockingChannel(b)
	runtime.GOMAXPROCS(1)
}

func benchmarkBlockingChannel(b *testing.B) {
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	channel := make(chan int64, 1024*16)
	go func() {
		for i := int64(0); i < iterations; i++ {
			channel <- i
		}
	}()

	for i := int64(0); i < iterations; i++ {
		msg := <-channel
		if msg != i {
			panic("Out of sequence")
		}
	}
}
