package main

func (this Worker) Process() int64 {
	next := this.sequence.Load() + 1
	available := this.barrier.Load()

	if next <= available {
		for next <= available {
			this.handler.Consume(next, available-next)
			next++
		}

		next--
		this.sequence.Store(next)
		return next
	} else if next <= this.source.Load() {
		return Gating
	} else {
		return Idle
	}
}

func NewWorker(barrier Barrier, handler Consumer, source, sequence *Sequence) Worker {
	// TODO: make this a pointer and test performance...
	return Worker{
		barrier:  barrier,
		handler:  handler,
		source:   source,
		sequence: sequence,
	}
}

type Worker struct {
	barrier  Barrier
	handler  Consumer
	source   *Sequence
	sequence *Sequence
}

type Consumer interface {
	Consume(sequence, remaining int64)
}

const (
	Gating int64 = -2
	Idle         = -3
)
