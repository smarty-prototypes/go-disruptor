package disruptor

import "testing"

func BenchmarkBarrierLoadSingle(b *testing.B) {
	barrier := NewBarrier(NewCursor())

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		barrier.Load()
	}
}

func BenchmarkBarrierLoadMultiple(b *testing.B) {
	barrier := NewBarrier(NewCursor(), NewCursor(), NewCursor(), NewCursor())

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		barrier.Load()
	}
}
