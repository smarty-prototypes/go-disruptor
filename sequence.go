package main

type Sequence [FillCPUCacheLine]int64

func (this *Sequence) Store(value int64) {
	(*this)[SequencePayloadIndex] = value
}
func (this *Sequence) Load() int64 {
	return (*this)[SequencePayloadIndex]
}

func NewSequence() *Sequence {
	return &Sequence{InitialSequenceValue}
}

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	SequencePayloadIndex uint8 = 0
	FillCPUCacheLine     uint8 = 8 // FUTURE: use build tags for i386, amd64, and ARM-v4,5,6,7,8 processors
)
