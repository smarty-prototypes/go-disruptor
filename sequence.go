package disruptor

type Sequence [FillCPUCacheLine]int64

func NewSequence() *Sequence {
	return &Sequence{InitialSequenceValue}
}

// TODO: aligned read/write: https://groups.google.com/forum/#!topic/golang-nuts/XDfQUn4U_g8

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	SequencePayloadIndex uint8 = 0
	FillCPUCacheLine     uint8 = 8
)
