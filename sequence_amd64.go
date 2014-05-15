package main

func (this *Sequence) Store(value int64) {
	(*this)[SequencePayloadIndex] = value
}
func (this *Sequence) Load() int64 {
	// atomic load runs FASTER? on amd64 (see comment from previous commit)
	//return atomic.LoadInt64(&(*this)[SequencePayloadIndex])
	return (*this)[SequencePayloadIndex]
}

const FillCPUCacheLine uint8 = 8
