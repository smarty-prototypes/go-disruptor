package main

func (this Worker) Process() uint8 {
	nextSequence := this.sequence.Load() + 1
	availableSequence := this.barrier.Load()

	if nextSequence <= availableSequence {
		for nextSequence <= availableSequence {
			this.handler.Consume(nextSequence, availableSequence-nextSequence)
			nextSequence++
		}

		this.sequence.Store(nextSequence - 1)
		return Processing
	} else if nextSequence <= this.source.Load() {
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
