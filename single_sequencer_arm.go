package disruptor

func (this *SingleProducerSequencer) Publish(sequence int64) {
	this.cursor.Store(sequence)
}
