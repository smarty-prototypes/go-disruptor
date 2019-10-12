package disruptor

import "sync/atomic"

type DefaultReader struct {
	state    int64
	current  *Cursor // this reader has processed up to this sequence
	written  *Cursor // the ring buffer has been written up to this sequence
	upstream Barrier // all of the readers have advanced up to this sequence
	waiter   WaitStrategy
	consumer Consumer
}

func NewReader(current, written *Cursor, upstream Barrier, waiter WaitStrategy, consumer Consumer) *DefaultReader {
	return &DefaultReader{
		state:    stateRunning,
		current:  current,
		written:  written,
		upstream: upstream,
		waiter:   waiter,
		consumer: consumer,
	}
}

func (this *DefaultReader) Listen() {
	var gateCount, idleCount, lower int64
	var upper = this.current.Load()

	for {
		lower = upper + 1
		upper = this.upstream.Load()

		if lower <= upper {
			this.consumer.Consume(lower, upper)
			this.current.Store(upper)
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
}

func (this *DefaultReader) Close() error {
	atomic.StoreInt64(&this.state, stateClosed)
	return nil
}

const (
	stateRunning = iota
	stateClosed
)
