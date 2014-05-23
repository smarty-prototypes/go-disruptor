package disruptor

type Reader struct {
	upstream Barrier
	written  *Cursor
	read     *Cursor
}

func NewReader(upstream Barrier, written, read *Cursor) *Reader {
	return &Reader{
		upstream: upstream,
		written:  written,
		read:     read,
	}
}

func (this *Reader) Receive() (int64, int64) {
	lower := this.read.Load() + 1
	upper := this.upstream.LoadBarrier(lower)

	if lower <= upper {
		return lower, upper
	} else if gate := this.written.Load(); lower <= gate {
		return InitialSequenceValue, Gating
	} else {
		return InitialSequenceValue, Idling
	}
}
