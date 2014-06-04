package benchmarks

import (
	"runtime"
	"testing"
)

func BenchmarkNonBlockingOneGoroutine(b *testing.B) {
	runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(1)
	benchmarkNonBlocking(b, 1)
}

func BenchmarkNonBlockingTwoGoroutines(b *testing.B) {
	runtime.GOMAXPROCS(2)
	defer runtime.GOMAXPROCS(1)
	benchmarkNonBlocking(b, 1)
}
func BenchmarkNonBlockingThreeGoroutinesWithContendedWrite(b *testing.B) {
	runtime.GOMAXPROCS(3)
	defer runtime.GOMAXPROCS(1)
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

	b.StopTimer()
}
