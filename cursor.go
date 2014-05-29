package disruptor

const (
	InitialSequenceValue int64 = -1
	cpuCacheLinePadding        = 7
)

const (
	Gating  = -2
	Idling  = -3
	Stopped = -4
)

type Cursor struct {
	sequence int64
	padding  [cpuCacheLinePadding]int64
}

func NewCursor() *Cursor {
	return &Cursor{sequence: InitialSequenceValue}
}

func (this *Cursor) Read(minimum int64) int64 {
	return this.sequence
}
