package disruptor

import "time"

type Reader struct {
	closed   *Cursor // indicates if the reader should continue processing
	read     *Cursor // this particular reader has advanced to this sequence
	written  *Cursor // the ring buffer has been written up to this sequence
	upstream Barrier // don't allow the reader to advance beyond this sequence
	consumer Consumer
}

func NewReader(read, written *Cursor, upstream Barrier, consumer Consumer) *Reader {
	return &Reader{
		closed:   NewCursor(),
		read:     read,
		written:  written,
		upstream: upstream,
		consumer: consumer,
	}
}

func (this *Reader) Listen() {
	previous := this.read.Load()
	idling, gating := 0, 0

	for {
		lower := previous + 1
		upper := this.upstream.Load()

		if lower <= upper {
			this.consumer.Consume(lower, upper)
			this.read.Store(upper)
			previous = upper
		} else if upper = this.written.Load(); lower <= upper {
			time.Sleep(time.Microsecond)
			// Gating--TODO: wait strategy (provide gating count to wait strategy for phased backoff)
			gating++
			idling = 0
		} else if this.closed.Load() == 0 {
			time.Sleep(time.Millisecond)
			// Idling--TODO: wait strategy (provide idling count to wait strategy for phased backoff)
			idling++
			gating = 0
		} else {
			break
		}
	}
}

func (this *Reader) Close() error {
	this.closed.Store(1)
	return nil
}
