package disruptor

func (this *Cursor) Store(value int64) {
	this.value = value
}
func (this *Cursor) Load() int64 {
	return this.value
}
