package disruptor

import "sync/atomic"

type DefaultListener struct {
	state    int64
	current  *Cursor // the listener has processed up to this sequence
	written  *Cursor // the ring buffer has been written up to this sequence
	upstream Barrier // the upstream listeners have processed up to this sequence
	waiter   WaitStrategy
	consumer Consumer
}

func NewListener(current, written *Cursor, upstream Barrier, waiter WaitStrategy, consumer Consumer) *DefaultListener {
	return &DefaultListener{
		state:    stateRunning,
		current:  current,
		written:  written,
		upstream: upstream,
		waiter:   waiter,
		consumer: consumer,
	}
}

func (this *DefaultListener) Listen() {
	var gateCount, idleCount, lower, upper int64
	var current = this.current.Load()

	for {
		lower = current + 1
		upper = this.upstream.Load()

		if lower <= upper {
			this.consumer.Consume(lower, upper)
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
}

func (this *DefaultListener) Close() error {
	atomic.StoreInt64(&this.state, stateClosed)
	return nil
}

const (
	stateRunning = iota
	stateClosed
)
