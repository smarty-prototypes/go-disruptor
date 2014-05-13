package main

import "fmt"

func (this Worker) Process() uint8 {
	current := this.sequence.Load()
	available := this.barrier.Load()

	if current+1 <= available {
		for ; current <= available; current++ {
			if current%1000000 == 0 {
				fmt.Println(current)
			}

			//this.handler.Consume(current, available-current)
		}

		this.sequence.Store(available + 1)
		return Processing
	} else if current+1 <= this.source.Load() {
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
