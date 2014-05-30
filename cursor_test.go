package disruptor

import "testing"

func BenchmarkWrites(b *testing.B) {
	iterations := int64(b.N)

	cursor := NewCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		cursor.Store(i)
	}
}
