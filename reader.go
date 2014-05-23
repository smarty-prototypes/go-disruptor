package disruptor

type Reader struct {
	read     *Cursor
	written  *Cursor
	upstream Barrier
}

func NewReader(read, written *Cursor, upstream Barrier) *Reader {
	return &Reader{
		read:     read,
		written:  written,
		upstream: upstream,
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
