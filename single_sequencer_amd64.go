package main

func (this *SingleProducerSequencer) Publish(sequence int64) {
	this.cursor[SequencePayloadIndex] = sequence
}
