package disruptor

import "testing"

func BenchmarkWriterReserve(b *testing.B) {
	iterations := int64(b.N)
	written, read := NewCursor(), NewCursor()
	writer := NewWriter(written, read, 1024*64)

	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		writer.Reserve(1)
		read.sequence = i
	}
}
