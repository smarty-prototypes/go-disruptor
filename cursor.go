package disruptor

const (
	InitialSequenceValue int64 = -1
	cpuCacheLinePadding        = 7
)

type Cursor struct {
	sequence int64
	padding  [cpuCacheLinePadding]int64
}

func NewCursor() *Cursor {
	return &Cursor{sequence: InitialSequenceValue}
}
