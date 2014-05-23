package disruptor

import "fmt"

const (
	Gating = -2
	Idling = -3
)

type Reader struct {
	upstreamBarrier Barrier
	writerCursor    *Cursor
	readerCursor    *Cursor
}

func NewReader(upstreamBarrier Barrier, writerCursor, readerCursor *Cursor) *Reader {
	return &Reader{
		upstreamBarrier: upstreamBarrier,
		writerCursor:    writerCursor,
		readerCursor:    readerCursor,
	}
}

// TODO: performance when current (or next?) sequence is received as a parameter to Receive
// instead of reading the cursor...
func (this *Reader) Receive() (int64, int64) {
	next := this.readerCursor.Load() + 1
	ready := this.upstreamBarrier.Load()
	fmt.Printf("\t\t\t\t\t[READER] Next: %d, Ready: %d\n", next, ready)

	if next <= ready {
		fmt.Printf("\t\t\t\t\t[READER] Next Sequence: %d, Remaining: %d\n", next, ready-next)
		return next, ready - next
	} else if gate := this.writerCursor.Load(); next <= gate {
		fmt.Println("\t\t\t\t\t[READER] Gating at sequence:", gate)
		return next, Gating
	} else {
		fmt.Println("\t\t\t\t\t[READER] Gating")
		return next, Idling
	}
}
