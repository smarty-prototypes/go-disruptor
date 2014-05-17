package disruptor

// TODO: rename to Commit
func (this *Writer) Publish(sequence int64) {
	this.writerCursor.value = sequence
}
