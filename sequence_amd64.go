package disruptor

func (this *Sequence) Store(value int64) {
	this[SequencePayloadIndex] = value
}
func (this *Sequence) Load() int64 {
	return this[SequencePayloadIndex]
}
