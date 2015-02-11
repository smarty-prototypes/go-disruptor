package disruptor

import "time"

type Reader struct {
	read     *Cursor
	written  *Cursor
	upstream Barrier
	consumer Consumer
	ready    bool
	done     chan struct{}
}

func NewReader(read, written *Cursor, upstream Barrier, consumer Consumer) *Reader {
	// start it closed, so calling Wait on an unstarted reader will return
	// instantly.
	done := make(chan struct{})
	close(done)

	return &Reader{
		read:     read,
		written:  written,
		upstream: upstream,
		consumer: consumer,
		ready:    false,
		done:     done,
	}
}

func (this *Reader) Start() {
	this.ready = true
	this.done = make(chan struct{})
	go this.receive(this.done)
}

func (this *Reader) StopAndWait() {
	this.Stop()
	this.Wait()
}

func (this *Reader) Stop() {
	this.ready = false
}

func (this *Reader) Wait() {
	<-this.done
}

func (this *Reader) receive(done chan struct{}) {
	previous := this.read.Load()
	idling, gating := 0, 0

	for {
		lower := previous + 1
		upper := this.upstream.Read(lower)

		if lower <= upper {
			this.consumer.Consume(lower, upper)
			this.read.Store(upper)
			previous = upper
		} else if upper = this.written.Load(); lower <= upper {
			// Gating--TODO: wait strategy (provide gating count to wait strategy for phased backoff)
			gating++
			idling = 0
		} else if this.ready {
			// Idling--TODO: wait strategy (provide idling count to wait strategy for phased backoff)
			idling++
			gating = 0
		} else {
			break
		}

		// sleeping increases the batch size which reduces number of writes required to store the sequence
		// reducing the number of writes allows the CPU to optimize the pipeline without prediction failures
		time.Sleep(time.Microsecond)
	}
	close(done)
}
