package disruptor

import "sync/atomic"

type defaultListener struct {
	state     *atomic.Int64
	current   *atomic.Int64   // processed up to this sequence
	committed sequenceBarrier // all writers have committed up to this sequence
	upstream  sequenceBarrier // upstream readers have completed up to this sequence
	waiting   WaitStrategy
	handler   Handler
}

func newListener(current *atomic.Int64, committed, upstream sequenceBarrier, waiting WaitStrategy, handler Handler) ListenCloser {
	return &defaultListener{
		state:     newAtomicInt64(stateRunning),
		current:   current,
		committed: committed,
		upstream:  upstream,
		waiting:   waiting,
		handler:   handler,
	}
}

func (this *defaultListener) Listen() {
	var gatedCount, idlingCount, lower, upper int64
	var current = this.current.Load()

	for {
		lower = current + 1
		upper = this.upstream.Load()

		if lower <= upper {
			this.handler.Handle(lower, upper)
			this.current.Store(upper)
			current = upper
		} else if upper = this.committed.Load(); lower <= upper {
			gatedCount++
			idlingCount = 0
			this.waiting.Gate(gatedCount)
		} else if this.state.Load() == stateRunning {
			idlingCount++
			gatedCount = 0
			this.waiting.Idle(idlingCount)
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
