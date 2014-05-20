package disruptor

import "testing"

func BenchmarkReader(b *testing.B) {
	writerCursor := NewCursor()
	readerCursor := NewCursor()
	writerBarrier := NewBarrier(writerCursor)

	reader := NewReader(writerBarrier, writerCursor, readerCursor)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	writerCursor.Store(1)

	for i := int64(0); i < iterations; i++ {
		sequence, _ := reader.Receive()
		reader.Commit(sequence)
		readerCursor.Store(0)
	}
}
