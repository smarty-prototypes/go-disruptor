package disruptor

import (
	"sync/atomic"
	"time"
)

type Reader struct {
	read     *Cursor // this particular reader has advanced to this sequence
	written  *Cursor // the ring buffer has been written up to this sequence
	upstream Barrier // don't allow the reader to advance beyond this sequence
	consumer Consumer
	ready    int32
}

func NewReader(read, written *Cursor, upstream Barrier, consumer Consumer) *Reader {
	return &Reader{
		read:     read,
		written:  written,
		upstream: upstream,
		consumer: consumer,
	}
}

func (this *Reader) Start() {
	atomic.StoreInt32(&this.ready, 0)
	go this.receive()
}
func (this *Reader) Stop() {
	atomic.StoreInt32(&this.ready, 1)
}
func (this *Reader) isReady() bool {
	return atomic.LoadInt32(&this.ready) == 0
}

func (this *Reader) receive() {
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
		} else if this.isReady() {
			time.Sleep(time.Millisecond)
			// Idling--TODO: wait strategy (provide idling count to wait strategy for phased backoff)
			idling++
			gating = 0
		} else {
			break
		}
	}
}
