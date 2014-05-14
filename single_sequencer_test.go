package main

import "testing"

func BenchmarkSingleProducerSequencerPublish(b *testing.B) {
	sequencer := NewSingleProducerSequencer(NewSequence(), 1024, Barrier{})
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequencer.Publish(i)
	}
}
