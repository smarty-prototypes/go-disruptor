package disruptor

func (this *Reader) Commit(sequence int64) {
	this.read.sequence = sequence
}
