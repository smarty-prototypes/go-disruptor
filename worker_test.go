package main

import "testing"

func BenchmarkWorker(b *testing.B) {
	producerSequence := NewSequence()
	workerSequence := NewSequence()
	barrier := NewBarrier(producerSequence)

	worker := NewWorker(barrier, &testHandler{}, producerSequence, workerSequence)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	producerSequence.Store(1)

	for i := int64(0); i < iterations; i++ {
		workerSequence.Store(0)
		worker.Process()
	}
}

type testHandler struct{}

func (this testHandler) Consume(sequence, remaining int64) {}
