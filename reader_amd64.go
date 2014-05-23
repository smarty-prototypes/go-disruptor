package disruptor

func (this *Reader) Commit(lower, upper int64) {
	this.read.sequence = upper
}
