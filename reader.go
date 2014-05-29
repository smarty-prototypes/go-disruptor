package disruptor

type Reader struct {
	read     *Cursor
	written  *Cursor
	upstream Barrier
	ready    bool
} // TODO: padding???

func NewReader(read, written *Cursor, upstream Barrier) *Reader {
	return &Reader{
		read:     read,
		written:  written,
		upstream: upstream,
		ready:    true,
	}
}

func (this *Reader) Start() {
	this.ready = true
}
func (this *Reader) Stop() {
	this.ready = false
}
func (this *Reader) Receive(next int64) int64 {
	maximum := this.upstream.Read(next)

	if next <= maximum {
		return maximum
	} else if maximum = this.written.Load(); next <= maximum {
		return Gating
	} else if this.ready {
		return Idling
	} else {
		return Stopped
	}
}

func (this *Reader) Commit(sequence int64) {
	this.read.Store(sequence)
}
