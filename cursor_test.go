package disruptor

import "testing"

func BenchmarkCursorStore(b *testing.B) {
	iterations := int64(b.N)

	cursor := NewCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		cursor.Store(i)
	}
}
func BenchmarkCursorLoad(b *testing.B) {
	iterations := int64(b.N)

	cursor := NewCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		cursor.Load()
	}
}
func BenchmarkCursorRead(b *testing.B) {
	iterations := int64(b.N)

	cursor := NewCursor()

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		cursor.Read(i)
	}
}
