package disruptor

import "time"

type Reader struct {
	read     *Cursor
	written  *Cursor
	upstream Barrier
	consumer Consumer
	ready    bool
} // TODO: padding???

func NewReader(read, written *Cursor, upstream Barrier, consumer Consumer) *Reader {
	return &Reader{
		read:     read,
		written:  written,
		upstream: upstream,
		consumer: consumer,
		ready:    false,
	}
}

func (this *Reader) Start() {
	this.ready = true
	go this.receive()
}
func (this *Reader) Stop() {
	this.ready = false
}

func (this *Reader) receive() {
	current := this.read.Sequence + 1
	for {
		gate := this.upstream.Read(current)

		if current <= gate {
			for current < gate {
				current += this.consumer.Consume(current, gate)
			}
			this.read.Store(current)
			current++
		} else if gate = this.written.Load(); current <= gate {
			// Gating--TODO: wait strategy (provide gating count to wait strategy for phased backoff)
			// gating++
			// idling = 0
			time.Sleep(time.Microsecond)
		} else if this.ready {
			// Idling--TODO: wait strategy (provide idling count to wait strategy for phased backoff)
			// idling++
			// gating = 0
			time.Sleep(time.Microsecond)
		} else {
			break
		}
	}
}
