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

func (this *Reader) Receive(lower int64) int64 {
	upper := this.upstream.LoadBarrier(lower)

	if lower <= upper {
		return upper
	} else if gate := this.written.Load(); lower <= gate {
		return Gating
	} else {
		return Idling
	}
}
