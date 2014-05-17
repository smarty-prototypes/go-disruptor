package disruptor

import "testing"

func BenchmarkCursorLoad(b *testing.B) {
	cursor := NewCursor()
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		cursor.Load()
	}
}

func BenchmarkCursorStore(b *testing.B) {
	cursor := NewCursor()
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		cursor.Store(i)
	}
}
