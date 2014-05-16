package disruptor

func (this *Sequence) Store(value int64) {
	this.value = value
}
func (this *Sequence) Load() int64 {
	return this.value
}
