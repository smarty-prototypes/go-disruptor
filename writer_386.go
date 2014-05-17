package disruptor

func (this *Writer) Commit(sequence int64) {
	this.writerCursor.Store(sequence)
}
