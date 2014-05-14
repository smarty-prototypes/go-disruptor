package main

type Sequence [FillCPUCacheLine]int64

func (this *Sequence) Store(value int64) {
	(*this)[0] = value
}
func (this *Sequence) Load() int64 {
	return (*this)[0]
}

func NewSequence() *Sequence {
	this := &Sequence{}
	this.Store(-1)
	return this
}

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	FillCPUCacheLine     uint8 = 8 // FUTURE: use build tags for i386, amd64, and ARM-v4,5,6,7,8 processors
)
