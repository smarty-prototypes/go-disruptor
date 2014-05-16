package disruptor

type Sequence struct {
	value   int64 // TODO: aligned read/write: https://groups.google.com/forum/#!topic/golang-nuts/XDfQUn4U_g8
	padding [cpuCacheLinePadding]int64
}

func NewSequence() *Sequence {
	return &Sequence{value: InitialSequenceValue}
}

const (
	MaxSequenceValue     int64 = (1 << 63) - 1
	InitialSequenceValue int64 = -1
	cpuCacheLinePadding  uint8 = 7
)
