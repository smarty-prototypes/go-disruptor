package disruptor

func (this *Reader) Commit(sequence int64) {
	this.readerCursor.Store(sequence)
}
