package disruptor

import "testing"

func BenchmarkSequenceLoad(b *testing.B) {
	sequence := NewSequence()
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence.Load()
	}
}

func BenchmarkSequenceStore(b *testing.B) {
	sequence := NewSequence()
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequence.Store(i)
	}
}
