package main

func (this Worker) Process() uint8 {
	next := this.sequence.Load() + 1
	available := this.barrier.Load()

	if next <= available {
		for next <= available {
			this.handler.Consume(next, available-next)
			next++
		}

		this.sequence.Store(next - 1)
		return Processing
	} else if next <= this.source.Load() {
		return Gating
	} else {
		return Idle
	}
}

func NewWorker(barrier Barrier, handler Consumer, source, sequence *Sequence) Worker {
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
	Processing uint8 = iota
	Gating
	Idle
)
