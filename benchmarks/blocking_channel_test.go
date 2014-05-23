package benchmarks

import "testing"

const blockingChannelBufferSize = 1024 * 1024

func BenchmarkBlockingChannel(b *testing.B) {
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	channel := make(chan int64, blockingChannelBufferSize)
	go func() {
		for i := int64(0); i < iterations; i++ {
			channel <- i
		}
	}()

	for i := int64(0); i < iterations; i++ {
		<-channel
	}
}
