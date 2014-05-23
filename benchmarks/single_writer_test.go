package benchmarks

import (
	"runtime"
	"testing"

	"github.com/smartystreets/go-disruptor"
)

const singleWriterRingBufferSize = 1024 * 16

func BenchmarkSingleWriter(b *testing.B) {
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	reader := disruptor.NewReader(read, written, written)
	writer := disruptor.NewWriter(written, read, singleWriterRingBufferSize)

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	go func() {
		for i := int64(0); i < iterations; i++ {
			if lower, upper := writer.Reserve(1); upper >= 0 {
				writer.Commit(lower, upper)
			}
		}
	}()

	for {
		_, upper := reader.Receive()
		if upper >= 0 {
			reader.Commit(upper)

			if upper+1 == iterations {
				break
			}
		}
	}
}

func init() {
	runtime.GOMAXPROCS(2)
}
