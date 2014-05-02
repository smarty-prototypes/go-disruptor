package disruptor

import "testing"

func BenchmarkChannels(b *testing.B) {
	iterations := int(b.N)
	channel := make(chan int, BufferSize)

	go func() {
		for sequence := 0; sequence < iterations; sequence++ {
			channel <- sequence
		}
	}()

	for sequence := 0; sequence < iterations; sequence++ {
		message := <-channel
		if message != sequence {
			panic("Out of sequence")
		}
	}
}
