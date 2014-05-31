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

// TODO: ARM, i386-specific methods
func (this *Cursor) Read(minimum int64) int64 {
	return this.sequence
}

func (this *Cursor) Load() int64 {
	return this.sequence
}
func (this *Cursor) Store(sequence int64) {
	this.sequence = sequence
}
