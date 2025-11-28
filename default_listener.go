package disruptor

import (
	"io"
	"sync/atomic"
)

type defaultListener struct {
	state    *atomic.Int64   // using atomic.Int64 which is padded to avoid false sharing
	current  *atomic.Int64   // this reader has processed up to this sequence
	written  *atomic.Int64   // the ring buffer has been written up to this sequence
	upstream sequenceBarrier // all readers have advanced to this sequence
	wait     WaitStrategy
	consumer Handler
}

func newListener(current, written *atomic.Int64, upstream sequenceBarrier, wait WaitStrategy, consumer Handler) ListenCloser {
	state := &atomic.Int64{}
	state.Store(stateRunning)
	return &defaultListener{
		state:    state,
		current:  current,
		written:  written,
		upstream: upstream,
		wait:     wait,
		consumer: consumer,
	}
}

func (this *defaultListener) Listen() {
	var gatedCount, idlingCount, lower, upper int64
	var current = this.current.Load()

	for {
		lower = current + 1
		upper = this.upstream.Load()

		if lower <= upper {
			this.consumer.Handle(lower, upper)
			this.current.Store(upper)
			current = upper
		} else if upper = this.written.Load(); lower <= upper {
			gatedCount++
			idlingCount = 0
			this.wait.Gate(gatedCount)
		} else if this.state.Load() == stateRunning {
			idlingCount++
			gatedCount = 0
			this.wait.Idle(idlingCount)
		} else {
			break
		}
	}

	if closer, ok := this.consumer.(io.Closer); ok {
		_ = closer.Close()
	}
}

func (this *defaultListener) Close() error {
	this.state.Store(stateClosed)
	return nil
}

const (
	stateRunning = iota
	stateClosed
)
