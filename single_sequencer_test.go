package main

import "testing"

func BenchmarkSingleProducerSequencerNext(b *testing.B) {
	consumerSequence := NewSequence()
	publisherSequence := NewSequence()
	consumerBarrier := NewBarrier(consumerSequence)

	sequencer := NewSingleProducerSequencer(publisherSequence, 1024, consumerBarrier)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		claimed := sequencer.Next(1)
		//consumerSequence[0] = claimed
		consumerSequence.Store(claimed)
	}
}

func BenchmarkSingleProducerSequencerPublish(b *testing.B) {
	sequencer := NewSingleProducerSequencer(NewSequence(), 1024, Barrier{})
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	for i := int64(0); i < iterations; i++ {
		sequencer.Publish(i)
	}
}
