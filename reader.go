package disruptor

import "time"

type Reader struct {
	closed   *Cursor
	reads    *Cursor // the reader has read up to this sequence
	writes   *Cursor // the ring buffer has been written up to this sequence
	upstream Barrier // the workers just in front of this reader have completed up to this sequence
	consumer Consumer
}

func NewReader(reads, writes *Cursor, upstream Barrier, consumer Consumer) *Reader {
	return &Reader{
		closed:   NewCursor(),
		reads:    reads,
		writes:   writes,
		upstream: upstream,
		consumer: consumer,
	}
}

func (this *Reader) Listen() {
	current := this.reads.Load()
	idling, gating := 0, 0

	for {
		lower := current + 1
		upper := this.upstream.Load()

		if lower <= upper {
			this.consumer.Consume(lower, upper)
			this.reads.Store(upper)
			current = upper
		} else if upper = this.writes.Load(); lower <= upper {
			time.Sleep(time.Microsecond)
			gating++
			idling = 0
		} else if this.closed.Load() == InitialCursorSequenceValue {
			time.Sleep(time.Millisecond)
			idling++
			gating = 0
		} else {
			break
		}
	}
}

func (this *Reader) Close() error {
	this.closed.Store(InitialCursorSequenceValue + 1)
	return nil
}
