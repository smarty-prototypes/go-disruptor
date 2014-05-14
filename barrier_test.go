package main

import "testing"

func BenchmarkBarrierLoad(b *testing.B) {
	upstream := NewSequence()
	upstream.Store(42)
	barrier := NewBarrier(upstream)

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		barrier.Load()
	}
}
