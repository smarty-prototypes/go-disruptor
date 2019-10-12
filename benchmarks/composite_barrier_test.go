package benchmarks

import (
	"testing"

	"github.com/smartystreets-prototypes/go-disruptor"
)

func BenchmarkCompositeBarrierRead(b *testing.B) {
	iterations := int64(b.N)

	barrier := disruptor.NewCompositeBarrier(
		disruptor.NewCursor(), disruptor.NewCursor(), disruptor.NewCursor(), disruptor.NewCursor())

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		barrier.Load()
	}
}
