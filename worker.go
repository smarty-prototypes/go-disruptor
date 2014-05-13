package main

func (this Worker) Process() uint8 {
	currentSequence := this.sequence.Load()
	nextSequence := currentSequence + 1
	availableSequence := this.barrier.Load()

	if nextSequence <= availableSequence {
		processedSequence := currentSequence

		for nextSequence <= availableSequence {
			this.handler.Consume(currentSequence+1, availableSequence-(currentSequence+1))
			processedSequence = nextSequence
			nextSequence++
		}

		this.sequence.Store(processedSequence)
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
