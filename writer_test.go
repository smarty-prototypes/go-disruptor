package disruptor

import "testing"

func BenchmarkWriterCommit(b *testing.B) {
	writer := NewWriter(NewCursor(), 1024, nil)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		writer.Commit(i, i)
	}
}

func BenchmarkWriterReserve(b *testing.B) {
	read := NewCursor()
	written := NewCursor()

	writer := NewWriter(written, 1024, read)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		claimed, _ := writer.Reserve(1)
		read.Store(claimed)
	}
}

func BenchmarkWriterNextWrapPoint(b *testing.B) {
	read := NewCursor()
	written := NewCursor()

	writer := NewWriter(written, 1024, read)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	read.Store(MaxSequenceValue)
	for i := int64(0); i < iterations; i++ {
		writer.Reserve(1)
	}
}
