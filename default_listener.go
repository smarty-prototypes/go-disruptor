package disruptor

import "sync/atomic"

// The defaultListener is single threaded and designed to run in a single goroutine. It tracks which slots or events in
// the associated ring buffer have been read, processed, or handled in some manner.
type defaultListener struct {
	state     *atomic.Int64
	current   *atomicSequence // the configured handler has processed up to this sequence
	committed sequenceBarrier // all sequencers (writers) have committed up to this sequence
	upstream  sequenceBarrier // any other configured groups handlers prior to or upstream have processed up to this sequence
	waiter    WaitStrategy
	handler   Handler
}

func newListener(current *atomicSequence, committed, upstream sequenceBarrier, waiter WaitStrategy, handler Handler) ListenCloser {
	return &defaultListener{
		state:     &atomic.Int64{},
		current:   current,
		committed: committed,
		upstream:  upstream,
		waiter:    waiter,
		handler:   handler,
	}
}

func (this *defaultListener) Listen() {
	var gatedCount, idlingCount, lower, upper int64
	var current = this.current.Load()

	for {
		lower = current + 1
		upper = this.upstream.Load(lower)

		if lower <= upper {
			this.handler.Handle(lower, upper)
			this.current.Store(upper)
			current = upper
			gatedCount = 0
			idlingCount = 0
		} else if upper = this.committed.Load(lower); lower <= upper {
			gatedCount++
			idlingCount = 0
			this.waiter.Gate(gatedCount)
		} else if this.state.Load() == stateRunning {
			idlingCount++
			gatedCount = 0
			this.waiter.Idle(idlingCount)
		} else {
			break
		}
	}
}

func (this *defaultListener) Close() error {
	this.state.Store(stateClosed)
	return nil
}

const (
	stateRunning = 0
	stateClosed  = 1
)
