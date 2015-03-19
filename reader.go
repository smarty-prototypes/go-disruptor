package disruptor

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"
)

type Reader struct {
	ID       uint64
	read     *Cursor
	written  *Cursor
	upstream Barrier
	consumer Consumer
	ready    bool
}

func NewReader(read, written *Cursor, upstream Barrier, consumer Consumer) *Reader {
	return &Reader{
		ID:       atomic.AddUint64(&readerID, 1),
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
	previous := this.read.Load()
	idling, gating := 0, 0

	for {
		lower := previous + 1
		upper := this.upstream.Read(lower)

		if lower <= upper {
			this.consumer.Consume(lower, upper)
			fmt.Printf("Reader %d completed sequence (lower: %d and upper:%d)\n", this.ID, lower, upper)
			this.read.Store(upper)
			previous = upper
		} else if upper1 := this.written.Load(); lower <= upper1 {
			runtime.Gosched()
			time.Sleep(time.Second)
			fmt.Printf("Reader %d gating on upstream at sequence (reader completed up to: %d, desired: %d, upstream gate: %d, written by writer :%d)\n", this.ID, lower-1, lower, upper, upper1)
			// Gating--TODO: wait strategy (provide gating count to wait strategy for phased backoff)
			gating++
			idling = 0
			upper = upper1
		} else if this.ready {
			runtime.Gosched()
			time.Sleep(time.Second)
			fmt.Printf("Reader %d idling on writer at sequence (reader completed up to: %d, desired: %d and available/upstream gate:%d)\n", this.ID, lower-1, lower, upper)
			// Idling--TODO: wait strategy (provide idling count to wait strategy for phased backoff)
			idling++
			gating = 0
		} else {
			break
		}

		// sleeping increases the batch size which reduces number of writes required to store the sequence
		// reducing the number of writes allows the CPU to optimize the pipeline without prediction failures
		time.Sleep(time.Microsecond) // TODO: runtime.Gosched()?
	}
}

var readerID uint64
