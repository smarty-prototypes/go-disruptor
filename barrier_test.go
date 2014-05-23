package disruptor

import "testing"

func BenchmarkBarrierLoadSingle(b *testing.B) {
	barrier := NewCursor()

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		barrier.LoadBarrier(0)
	}
}
