package main

func (this Worker) Process(consumer int) uint8 {
	next := this.sequence.Load() + 1
	available := this.barrier.Load(consumer)

	if next <= available {
		// fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: %d work items found.\n", consumer, available-next+1)

		for next <= available {
			// fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Consuming sequence %d\n", consumer, next)

			this.handler.Consume(next, available-next)
			next++
		}

		// fmt.Printf("\t\t\t\t\t\t\t\t\tConsumer %d:: Completed through sequence %d\n", consumer, next-1)
		this.sequence.Store(next - 1)
		return Processing
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
	Processing uint8 = iota
	Gating
	Idle
)
