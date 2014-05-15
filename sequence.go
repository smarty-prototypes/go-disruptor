package disruptor

type Sequence [FillCPUCacheLine]int64

func NewSequence() *Sequence {
	return &Sequence{InitialSequenceValue}
}

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	SequencePayloadIndex uint8 = 0
)
