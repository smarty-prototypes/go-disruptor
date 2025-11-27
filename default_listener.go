package disruptor

import (
	"io"
	"sync/atomic"
)

type defaultListener struct {
	state    int64
	current  *atomic.Int64   // this reader has processed up to this sequence
	written  *atomic.Int64   // the ring buffer has been written up to this sequence
	upstream sequenceBarrier // all of the readers have advanced up to this sequence
	waiter   WaitStrategy
	consumer Handler
}

func newListener(current, written *atomic.Int64, upstream sequenceBarrier, waiter WaitStrategy, consumer Handler) ListenCloser {
	return &defaultListener{
		state:    stateRunning,
		current:  current,
		written:  written,
		upstream: upstream,
		waiter:   waiter,
		consumer: consumer,
	}
}

func (this *defaultListener) Listen() {
	var gateCount, idleCount, lower, upper int64
	var current = this.current.Load()

	for {
		lower = current + 1
		upper = this.upstream.Load()

		if lower <= upper {
			this.consumer.Handle(lower, upper)
			this.current.Store(upper)
			current = upper
		} else if upper = this.written.Load(); lower <= upper {
			gateCount++
			idleCount = 0
			this.waiter.Gate(gateCount)
		} else if atomic.LoadInt64(&this.state) == stateRunning {
			idleCount++
			gateCount = 0
			this.waiter.Idle(idleCount)
		} else {
			break
		}
	}

	if closer, ok := this.consumer.(io.Closer); ok {
		_ = closer.Close()
	}
}

func (this *defaultListener) Close() error {
	atomic.StoreInt64(&this.state, stateClosed)
	return nil
}

const (
	stateRunning = iota
	stateClosed
)
