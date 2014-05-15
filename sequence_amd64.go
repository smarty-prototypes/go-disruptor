package disruptor

func (this *Sequence) Store(value int64) {
	(*this)[SequencePayloadIndex] = value
}
func (this *Sequence) Load() int64 {
	return (*this)[SequencePayloadIndex] // scheduler causes atomic load to run faster?
}

const FillCPUCacheLine uint8 = 8
