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
	Sequence int64
	padding  [cpuCacheLinePadding]int64
}

func NewCursor() *Cursor {
	return &Cursor{Sequence: InitialSequenceValue}
}

func (this *Cursor) Read(minimum int64) int64 {
	return this.Sequence
}

func (this *Cursor) Load() int64 {
	return this.Sequence
}
func (this *Cursor) Store(sequence int64) {
	this.Sequence = sequence
}
