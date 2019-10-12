package benchmarks

import (
	"testing"

	"github.com/smartystreets-prototypes/go-disruptor"
)

func BenchmarkCursorStore(b *testing.B) {
	iterations := int64(b.N)
	sequence := disruptor.NewCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence.Store(i)
	}
}
func BenchmarkCursorLoad(b *testing.B) {
	iterations := int64(b.N)
	sequence := disruptor.NewCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		_ = sequence.Load()
	}
}

func BenchmarkCursorLoadAsBarrier(b *testing.B) {
	var barrier disruptor.Barrier = disruptor.NewCursor()
	iterations := int64(b.N)

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		_ = barrier.Load()
	}
}
