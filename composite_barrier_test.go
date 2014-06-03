package disruptor

import "testing"

func BenchmarkCompositeBarrierRead(b *testing.B) {
	barrier := NewCompositeBarrier(NewCursor(), NewCursor(), NewCursor(), NewCursor())

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		barrier.Read(0)
	}
}
