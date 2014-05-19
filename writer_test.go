package disruptor

import "testing"

func BenchmarkWriterCommit(b *testing.B) {
	writer := NewWriter(NewCursor(), 1024, nil)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		writer.Commit(i)
	}
}

func BenchmarkWriterReserve(b *testing.B) {
	readerCursor := NewCursor()
	writerCursor := NewCursor()
	readerBarrier := NewBarrier(readerCursor)

	writer := NewWriter(writerCursor, 1024, readerBarrier)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		claimed := writer.Reserve(1)
		readerCursor.Store(claimed)
	}
}

func BenchmarkWriterNextWrapPoint(b *testing.B) {
	readerCursor := NewCursor()
	writerCursor := NewCursor()
	readerBarrier := NewBarrier(readerCursor)

	writer := NewWriter(writerCursor, 1024, readerBarrier)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	readerCursor.Store(MaxSequenceValue)
	for i := int64(0); i < iterations; i++ {
		writer.Reserve(1)
	}
}
