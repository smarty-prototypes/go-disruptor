package benchmarks

import (
	"runtime"
	"testing"
)

func BenchmarkNonBlockingChannelOneGoroutine(b *testing.B) {
	benchmarkNonBlockingChannel(b)
}

func BenchmarkNonBlockingChannelTwoGoroutines(b *testing.B) {
	runtime.GOMAXPROCS(2)
	benchmarkNonBlockingChannel(b)
	runtime.GOMAXPROCS(1)
}

func benchmarkNonBlockingChannel(b *testing.B) {
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	channel := make(chan int64, 1024*16)
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

	for i := int64(0); i < iterations; {
		select {
		case msg := <-channel:
			if msg != i {
				panic("Out of sequence")
			} else {
				i++
			}
		default:
			continue
		}
	}
}
