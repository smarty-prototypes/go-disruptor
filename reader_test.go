package disruptor

import "testing"

func BenchmarkReader(b *testing.B) {
	writerCursor := NewCursor()
	readerCursor := NewCursor()
	writerBarrier := NewBarrier(writerCursor)

	reader := NewReader(writerBarrier, &testHandler{}, writerCursor, readerCursor)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	writerCursor.Store(1)

	for i := int64(0); i < iterations; i++ {
		readerCursor.Store(0)
		reader.Process()
	}
}

type testHandler struct{}

func (this testHandler) Consume(sequence, remaining int64) {}
