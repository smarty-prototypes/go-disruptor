package benchmarks

import (
	"math"
	"runtime"
	"testing"

	"github.com/smartystreets/go-disruptor"
)

const multiWriterRingBufferSize = 1024 * 16
const multiWriterRingBufferMask = multiWriterRingBufferSize - 1

func BenchmarkDisruptorSharedWriterClaimSingle(b *testing.B) {
	benchmarkMultiWriter(b, 1)
}
func BenchmarkDisruptorSharedWriterClaimMultiple(b *testing.B) {
	benchmarkMultiWriter(b, 32)
}

func benchmarkMultiWriter(b *testing.B, maxClaim int64) {
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	shared := disruptor.NewSharedWriterBarrier(written, multiWriterRingBufferSize)
	reader := disruptor.NewReader(read, written, shared)
	writer := disruptor.NewSharedWriter(shared, read)

	iterations := int64(b.N)
	ringBuffer := [multiWriterRingBufferSize]int64{}
	claim := (int64(math.Log2(float64(iterations))) + 1)
	if claim >= maxClaim {
		claim = maxClaim
	}

	b.ReportAllocs()
	b.ResetTimer()

	go func() {
		sequence := int64(0)
		for sequence < iterations {
			if lower, upper := writer.Reserve(claim); upper >= 0 {
				for ; sequence <= upper; sequence++ {
					ringBuffer[sequence&multiWriterRingBufferMask] = sequence
				}

				writer.Commit(lower, upper)
				sequence = upper
			}
		}
	}()

	sequence := int64(0)
	for sequence < iterations {
		if _, upper := reader.Receive(); upper >= 0 {
			for ; sequence <= upper; sequence++ {
				if sequence != ringBuffer[sequence&multiWriterRingBufferMask] {
					panic("Out of sequence")
				}
			}

			reader.Commit(upper)
			sequence = upper
		}
	}
}

func init() {
	runtime.GOMAXPROCS(2)
}
