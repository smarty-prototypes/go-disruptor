package disruptor

import "testing"

func BenchmarkReader(b *testing.B) {
	written := NewCursor()
	read := NewCursor()
	reader := NewReader(written, written, read)

	written.Store(1)

	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		_, upper := reader.Receive()
		reader.Commit(upper)
		read.Store(0)
	}
}
