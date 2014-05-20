package benchmarks

import "testing"

const nonBlockingChannelBufferSize = 1024 * 1024

func BenchmarkNonBlockingChannel(b *testing.B) {
	iterations := int64(b.N)

	channel := make(chan int64, nonBlockingChannelBufferSize)
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

	for i := int64(0); i < iterations; i++ {
		select {
		case <-channel:
			i++
		default:
			continue
		}
	}
}
