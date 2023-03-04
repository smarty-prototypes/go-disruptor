package disruptor

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type DefaultWriter struct {
	mtx sync.Mutex

	written  *Cursor // the ring buffer has been written up to this sequence
	upstream Barrier // all of the readers have advanced up to this sequence
	capacity int64
	previous int64
	gate     int64
}

func NewWriter(written *Cursor, upstream Barrier, capacity int64) *DefaultWriter {
	return &DefaultWriter{
		upstream: upstream,
		written:  written,
		capacity: capacity,
		previous: defaultCursorValue,
		gate:     defaultCursorValue,
	}
}

func (this *DefaultWriter) Reserve(count int64) int64 {
	if count <= 0 {
		panic(ErrMinimumReservationSize)
	}

	// TODO 加了这个锁，保证多生产者并行发送无问题。
	//this.mtx.Lock()
	//defer this.mtx.Unlock()

	//this.previous += count
	//for spin := int64(0); this.previous-this.capacity > this.gate; spin++ {
	//	if spin&SpinMask == 0 {
	//		runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
	//	}
	//
	//	this.gate = this.upstream.Load()
	//}
	//return this.previous

	previous := atomic.AddInt64(&this.previous, count)
	for spin := int64(0); previous-this.capacity > this.gate; spin++ {
		if spin&SpinMask == 0 {
			runtime.Gosched() // LockSupport.parkNanos(1L); http://bit.ly/1xiDINZ
		}

		this.gate = this.upstream.Load()
	}

	return previous
}

func (this *DefaultWriter) Commit(_, upper int64) { this.written.Store(upper) }

const SpinMask = 1024*16 - 1 // arbitrary; we'll want to experiment with different values
