package benchmarks

import (
	"testing"

	"github.com/smartystreets-prototypes/go-disruptor"
)

func BenchmarkSequenceStore(b *testing.B) {
	iterations := int64(b.N)
	sequence := disruptor.NewSequence()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence.Store(i)
	}
}
func BenchmarkSequenceLoad(b *testing.B) {
	iterations := int64(b.N)
	sequence := disruptor.NewSequence()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		_ = sequence.Load()
	}
}

func BenchmarkSequenceLoadAsBarrier(b *testing.B) {
	var barrier disruptor.Barrier = disruptor.NewSequence()
	iterations := int64(b.N)

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		_ = barrier.Load()
	}
}
