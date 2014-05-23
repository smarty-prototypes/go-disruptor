package disruptor

func (this *Reader) Commit(upper int64) {
	this.read.Store(upper)
}
