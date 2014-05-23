package disruptor

func (this *Reader) Commit(lower, upper int64) {
	this.readerCursor.Store(upper)
}
