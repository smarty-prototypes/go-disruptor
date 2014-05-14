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
		consumerSequence.Store(claimed)
	}
}

func BenchmarkSingleProducerSequencerNextWrap(b *testing.B) {
	consumerSequence := NewSequence()
	publisherSequence := NewSequence()
	consumerBarrier := NewBarrier(consumerSequence)

	sequencer := NewSingleProducerSequencer(publisherSequence, 1024, consumerBarrier)
	iterations := int64(b.N)
	b.ReportAllocs()
	b.ResetTimer()

	consumerSequence.Store(MaxSequenceValue)
	for i := int64(0); i < iterations; i++ {
		sequencer.Next(1)
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
