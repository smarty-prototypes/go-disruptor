package benchmarks

import (
	"testing"
)

func BenchmarkBlockingOneGoroutine(b *testing.B) {
	benchmarkBlocking(b, 1)
}
func BenchmarkBlockingTwoGoroutines(b *testing.B) {
	benchmarkBlocking(b, 1)
}
func BenchmarkBlockingThreeGoroutinesWithContendedWrite(b *testing.B) {
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
			panic("out of sequence")
		}
	}

	b.StopTimer()
}

func BenchmarkNonBlockingOneGoroutine(b *testing.B) {
	benchmarkNonBlocking(b, 1)
}
func BenchmarkNonBlockingTwoGoroutines(b *testing.B) {
	benchmarkNonBlocking(b, 1)
}
func BenchmarkNonBlockingThreeGoroutinesWithContendedWrite(b *testing.B) {
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
