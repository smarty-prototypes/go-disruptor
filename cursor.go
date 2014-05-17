package disruptor

const (
	InitialCursorValue        = -1
	cpuCacheLinePadding       = 7
	MaxCursorValue      int64 = (1 << 63) - 1
)

type Cursor struct {
	value   int64 // TODO: aligned read/write: https://groups.google.com/forum/#!topic/golang-nuts/XDfQUn4U_g8
	padding [cpuCacheLinePadding]int64
}

func NewCursor() *Cursor {
	return &Cursor{value: InitialCursorValue}
}
