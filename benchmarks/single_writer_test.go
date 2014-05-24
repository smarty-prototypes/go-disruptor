package benchmarks

import (
	"math"
	"runtime"
	"testing"

	"github.com/smartystreets/go-disruptor"
)

const singleWriterRingBufferSize = 1024 * 16
const singleWriterRingBufferMask = singleWriterRingBufferSize - 1

func BenchmarkDisruptorWriterSingleClaim(b *testing.B) {
	benchmarkSingleWriter(b, 1)
}
func BenchmarkDisruptorWriterClaimMultiple(b *testing.B) {
	benchmarkSingleWriter(b, 32)
}

func benchmarkSingleWriter(b *testing.B, maxClaim int64) {
	written, read := disruptor.NewCursor(), disruptor.NewCursor()
	reader := disruptor.NewReader(read, written, written)
	writer := disruptor.NewWriter(written, read, singleWriterRingBufferSize)

	iterations := int64(b.N)
	ringBuffer := [singleWriterRingBufferSize]int64{}
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
					ringBuffer[sequence&singleWriterRingBufferMask] = sequence
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
				if sequence != ringBuffer[sequence&singleWriterRingBufferMask] {
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
